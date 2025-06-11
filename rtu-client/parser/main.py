import sqlite3
import xml.etree.ElementTree as ET


def get_ns(fp: str) -> dict[str, str]:
    ns = {}
    for event, elem in ET.iterparse(fp, events=("start-ns",)):
        prefix, uri = elem
        ns[prefix] = uri
    return ns


def main():
    fp = "scd/CEZ_VNlinka_LAB.scd"
    tree = ET.parse(fp)
    root = tree.getroot()
    events = ["start", "start-ns"]
    ns = get_ns(fp)
    conn = sqlite3.connect("first.db")
    c = conn.cursor()

    # 1) Substations / VoltageLevels / Bays
    substations = []
    voltage_levels = []
    bays = []

    for sub in root.findall("Substation", ns):
        sub_name = sub.get("name")
        sub_desc = sub.get("desc")
        # substations.append({"name": sub_name, "desc": sub.get("desc")})
        c.execute(
            "INSERT INTO substation(name, desc) VALUES(?, ?)",
            [sub_name, sub_desc],
        )
        c.execute("SELECT id FROM substation WHERE name = ?", [sub_name])
        sub_id = c.fetchone()[0]
        for vl in sub.findall("VoltageLevel", ns):
            vl_name = vl.get("name")
            vl_desc = vl.get("desc")
            c.execute(
                "INSERT INTO voltage_level(substation_id, name, desc) VALUES(?, ?, ?)",
                [sub_id, vl_name, vl_desc],
            )
            c.execute("SELECT id FROM voltage_level WHERE name = ?", [vl_name])
            vl_id = c.fetchone()[0]
            for bay in vl.findall("Bay", ns):
                bay_name = bay.get("name")
                bay_desc = bay.get("desc")
                c.execute(
                    "INSERT INTO bay(voltage_level_id, name, desc) VALUES(?, ?, ?)",
                    [vl_id, bay_name, bay_desc],
                )
    # 2) IED hierarchy
    ieds = []
    access_points = []
    servers = []
    ldevices = []
    ln0s = []
    datasets = []
    fcdas = []

    for ied in root.findall("IED", ns):
        ieds.append(
            {
                "name": ied.get("name"),
                "type": ied.get("type"),
                "manufacturer": ied.get("manufacturer"),
                "config_version": ied.get("configVersion"),
            }
        )
        for ap in ied.findall("AccessPoint", ns):
            access_points.append({"ied_name": ap.get("name"), "name": ap.get("name")})
            for srv in ap.findall("Server", ns):
                servers.append({"access_point_name": ap.get("name")})
                for ld in srv.findall("LDevice", ns):
                    ldevices.append(
                        {"access_point_name": ap.get("name"), "inst": ld.get("inst")}
                    )
                    # LN0
                    ln0 = ld.find("LN0", ns)
                    if ln0 is not None:
                        ln0s.append(
                            {
                                "ldevice_inst": ld.get("inst"),
                                "lnClass": ln0.get("lnClass"),
                                "inst": ln0.get("inst"),
                                "lnType": ln0.get("lnType"),
                            }
                        )
                        # datasets under LN0
                        for ds in ln0.findall("DataSet", ns):
                            datasets.append(
                                {
                                    "ln0_inst": ln0.get("inst"),
                                    "name": ds.get("name"),
                                    "desc": ds.get("desc"),
                                }
                            )
                            for fcda in ds.findall("FCDA", ns):
                                fcdas.append(
                                    {
                                        "dataset_name": ds.get("name"),
                                        "ldInst": fcda.get("ldInst"),
                                        "prefix": fcda.get("prefix"),
                                        "lnClass": fcda.get("lnClass"),
                                        "lnInst": fcda.get("lnInst"),
                                        "doName": fcda.get("doName"),
                                        "fc": fcda.get("fc"),
                                    }
                                )
    # print("ieds: ", ieds)
    # print("aps: ", access_points)
    # print("servers ", servers)
    # print("ldevices: ", ldevices)
    # print("ln0s: ", ln0s)
    # print("datasets: ", datasets)
    # print("fcdas: ", fcdas)
    #
    # sub_id_map = {name: _id for _id, name in c.fetchall()}
    # print("substations", substations)
    # print("sub_id_map", sub_id_map)
    # print("voltage_levels", voltage_levels)
    #
    # # 2) Insert voltage_levels (attach substation_id)
    # for vl in voltage_levels:
    #     vl["substation_id"] = sub_id_map[vl.pop("substation_name")]
    #     print(vl)
    # c.executemany(
    #     "INSERT INTO voltage_level(substation_id, name, desc) "
    #     "VALUES(:substation_id, :name, :desc)",
    #     voltage_levels,
    # )
    # # build (sub_name,vl_name)→id map
    # c.execute(
    #     "SELECT vl.id, s.name, vl.name FROM voltage_level vl JOIN substation s ON vl.substation_id=s.id"
    # )
    # vl_id_map = {(sub_name, vl_name): _id for _id, sub_name, vl_name in c.fetchall()}

    # # 3) Insert bays (attach voltage_level_id)
    # for bay in bays:
    #     key = (None, bay["voltage_level_name"])
    #     # find matching substation in voltage_levels list, if needed
    #     # here we assume voltage_level names are unique across substations
    #     bay["voltage_level_id"] = vl_id_map[key]
    #     bay.pop("voltage_level_name")
    # c.executemany(
    #     "INSERT INTO bay(voltage_level_id, name, desc) "
    #     "VALUES(:voltage_level_id, :name, :desc)",
    #     bays,
    # )
    #
    # # 4) Insert IEDs
    # c.executemany(
    #     "INSERT INTO ied(name, type, manufacturer, config_version) "
    #     "VALUES(:name, :type, :manufacturer, :config_version)",
    #     ieds,
    # )
    # c.execute("SELECT id, name FROM ied")
    # ied_id_map = {name: _id for _id, name in c.fetchall()}
    #
    # # 5) Insert AccessPoints
    # for ap in access_points:
    #     ap["ied_id"] = ied_id_map[ap.pop("ied_name")]
    # c.executemany(
    #     "INSERT INTO access_point(ied_id, name) " "VALUES(:ied_id, :name)",
    #     access_points,
    # )
    # c.execute("SELECT id, name FROM access_point")
    # ap_id_map = {name: _id for _id, name in c.fetchall()}
    #
    # # 6) Insert Servers
    # for srv in servers:
    #     srv["access_point_id"] = ap_id_map[srv.pop("access_point_name")]
    # c.executemany(
    #     "INSERT INTO server(access_point_id) " "VALUES(:access_point_id)", servers
    # )
    # c.execute("SELECT id, access_point_id FROM server")
    # srv_id_map = {ap_id: _id for _id, ap_id in c.fetchall()}
    #
    # # 7) Insert LDevices
    # for ld in ldevices:
    #     ld["server_id"] = srv_id_map[ap_id_map[ld.pop("server_ap_name")]]
    # c.executemany(
    #     "INSERT INTO ldevice(server_id, inst) " "VALUES(:server_id, :inst)", ldevices
    # )
    # c.execute("SELECT id, inst FROM ldevice")
    # ld_id_map = {inst: _id for _id, inst in c.fetchall()}
    #
    # # 8) Insert LN0s
    # for ln in ln0s:
    #     ln["ldevice_id"] = ld_id_map[ln.pop("ldevice_inst")]
    # c.executemany(
    #     "INSERT INTO ln0(ldevice_id, ln_class, inst, ln_type) "
    #     "VALUES(:ldevice_id, :ln_class, :inst, :ln_type)",
    #     ln0s,
    # )
    # c.execute("SELECT id, inst FROM ln0")
    # ln0_id_map = {inst: _id for _id, inst in c.fetchall()}
    #
    # # 9) Insert DataSets
    # for ds in datasets:
    #     ds["ln0_id"] = ln0_id_map[ds.pop("ln0_inst")]
    # c.executemany(
    #     "INSERT INTO dataset(ln0_id, name, desc) " "VALUES(:ln0_id, :name, :desc)",
    #     datasets,
    # )
    # c.execute("SELECT id, name FROM dataset")
    # ds_id_map = {name: _id for _id, name in c.fetchall()}
    #
    # # 10) Insert FCDAs
    # for f in fcdas:
    #     f["dataset_id"] = ds_id_map[f.pop("dataset_name")]
    # c.executemany(
    #     "INSERT INTO fcda(dataset_id, ld_inst, prefix, ln_class, ln_inst, do_name, fc) "
    #     "VALUES(:dataset_id, :ld_inst, :prefix, :ln_class, :ln_inst, :do_name, :fc)",
    #     fcdas,
    # )
    #
    # conn.commit()
    conn.commit()
    conn.close()

    # print("Data inserted successfully into scd.db")


if __name__ == "__main__":
    main()
