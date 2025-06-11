import json
import xml.etree.ElementTree as ET
from pprint import pprint

from utils import format_tag, print_rle


def parse_element(ele: ET.Element, indent: int = 0):
    prefix = " " * indent
    tag = format_tag(ele.tag)
    attribs = dict(ele.attrib.items())
    children = list(ele)

    print(f"{prefix} Tag: {tag}")
    if attribs:
        print(f"{prefix} Attributes:")
        attribs_json = json.dumps(attribs, indent=4)
        for line in attribs_json.splitlines():
            print(f"{prefix}    {line}")

    if children:
        child_tags = [format_tag(c.tag) for c in children]
        print(f"{prefix}  Children ({len(children)}): {child_tags}")
        for child in children:
            parse_element(child, indent + 4)


def get_ns(fp: str) -> dict[str, str]:
    ns = {}
    for event, elem in ET.iterparse(fp, events=("start-ns",)):
        prefix, uri = elem
        ns[prefix] = uri
    return ns


if __name__ == "__main__":
    events = ["start", "start-ns"]
    fp = "scd/CEZ_VNlinka_LAB.scd"
    ns = get_ns(fp)
    tree = ET.parse(fp)
    root = tree.getroot()
    tags = []
    for child in root:
        tag = format_tag(child.tag)
        tags.append(tag)
    # print_rle(tags, "Tag: {f}, count: {s}")

    for ele in root.findall("IED", ns):
        parse_element(ele)
