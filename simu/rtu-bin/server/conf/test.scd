<?xml version="1.0" encoding="UTF-8"?>
<SCL xmlns="http://www.iec.ch/61850/2003/SCL">
	<Header id=""/>
	<IED name="SampleIED">
		<Services>
			<DynAssociation />
			<GetDirectory />
			<GetDataObjectDefinition />
			<GetDataSetValue />
			<DataSetDirectory />
			<ReadWrite />
			<GetCBValues />
			<ConfLNs fixPrefix="true" fixLnInst="true" />
			<FileHandling />
			<TimerActivatedControl />
		</Services>
		<AccessPoint name="accessPoint1">
			<Server>
				<Authentication />
				<LDevice inst="Device1">
					<LN0 lnClass="LLN0" lnType="LLN01" inst="">

						<DataSet name="dataset1" desc="">
							<FCDA ldInst="Device1" lnClass="LLN0" fc="ST" doName="Mod" daName="q" />
							<FCDA ldInst="Device1" lnClass="MMXU" fc="ST" lnInst="1" doName="Mod" daName="q" />
							<FCDA ldInst="Device1" lnClass="MMXU" fc="CF" lnInst="1" doName="Mod" daName="ctlModel" />
						</DataSet>

						<ReportControl name="LLN0_Events_BuffRep" rptID="LLN0$RP$brcbEV1" datSet="dataset1" confRev="1" buffered="true"
							desc="Events brcb buffered report" intgPd="900000" bufTime="50">
						<TrgOps dchg="true" qchg="false" dupd="false" period="true" />
						<OptFields seqNum="true" timeStamp="true" dataSet="true" reasonCode="true" dataRef="false" entryID="true"
							configRef="true" bufOvfl="true" />
						<RptEnabled max="1" />
					</ReportControl>

					<DOI name="Mod">
						<DAI name="ctlModel">
							<Val>status-only</Val>
						</DAI>
					</DOI>
				</LN0>
				<LN lnClass="LPHD" lnType="LPHD1" inst="1" prefix="" />
				<LN lnClass="DGEN" lnType="DGEN1" inst="1" prefix="">
					<DOI name="Mod">
						<DAI name="ctlModel">
							<Val>status-only</Val>
						</DAI>
					</DOI>
				</LN>
				<LN lnClass="DSCH" lnType="DSCH1" inst="1" prefix="">
					<DOI name="Mod">
						<DAI name="ctlModel">
							<Val>status-only</Val>
						</DAI>
					</DOI>
				</LN>
				<LN lnClass="MMXU" lnType="MMXU1" inst="1" prefix="">
					<DOI name="Mod">
						<DAI name="ctlModel">
							<Val>status-only</Val>
						</DAI>
					</DOI>
				</LN>
				<LN lnClass="MMXU" lnType="MMXU2" inst="2" prefix="">
					<DOI name="Mod">
						<DAI name="ctlModel">
							<Val>status-only</Val>
						</DAI>
					</DOI>
				</LN>
			</LDevice>
		</Server>
	</AccessPoint>
</IED>
<DataTypeTemplates>
	<LNodeType id="LLN01" lnClass="LLN0">
		<DO name="Mod" type="INC_1_Mod" />
		<DO name="Beh" type="INS_1_Beh" />
		<DO name="Health" type="INS_1_Beh" />
		<DO name="NamPlt" type="LPL_1_NamPlt" />
	</LNodeType>
	<LNodeType id="LPHD1" lnClass="LPHD">
		<DO name="PhyNam" type="DPL_1_PhyNam" />
		<DO name="PhyHealth" type="INS_1_Beh" />
		<DO name="Proxy" type="SPS_1_Proxy" />
	</LNodeType>
	<LNodeType id="MMXU2" lnClass="MMXU" desc="Measurement">
		<DO name="Mod" type="INC_1_Mod" />
		<DO name="Beh" type="INS_1_Beh" />
		<DO name="Health" type="INS_1_Beh" />
		<DO name="NamPlt" type="LPL_2_NamPlt" />
		<DO name="TotW" type="MV_1_TotWh" />
	</LNodeType>
	<LNodeType id="DGEN1" lnClass="DGEN" desc="DER unit generator">
		<DO name="Mod" type="ENC_1_Mod" />
		<DO name="Beh" type="ENS_1_Beh" />
		<DO name="Health" type="ENS_1_Beh" />
		<DO name="NamPlt" type="LPL_2_NamPlt" />
		<DO name="OpTmh" type="INS_1_Beh" />
		<DO name="GnOpSt" type="ENS_1_Beh" />
		<DO name="OpTmsRs" type="INS_1_Beh" />
		<DO name="TotWh" type="MV_1_TotWh" />
	</LNodeType>
	<LNodeType id="DSCH1" lnClass="DSCH" desc="DER energy and/or ancillary services schedule">
		<DO name="Mod" type="ENC_1_Mod" />
		<DO name="Beh" type="ENS_1_Beh" />
		<DO name="Health" type="ENS_1_Beh" />
		<DO name="NamPlt" type="LPL_2_NamPlt" />
		<DO name="SchdSt" type="INS_1_Beh" />
		<DO name="SchdId" type="ING_1_SchdId" />
		<DO name="SchdCat" type="ING_1_SchdId" />
		<DO name="SchdTyp" type="ING_1_SchdId" />
		<DO name="SchdAbsTm" type="SCA_1_SchdAbsTm" />
	</LNodeType>
	<LNodeType id="MMXU1" lnClass="MMXU" desc="Measurement">
		<DO name="Mod" type="INC_1_Mod" />
		<DO name="Beh" type="INS_1_Beh" />
		<DO name="Health" type="INS_1_Beh" />
		<DO name="NamPlt" type="LPL_2_NamPlt" />
	</LNodeType>
	<DOType id="INC_1_Mod" cdc="INC">
		<DA name="q" bType="Quality" fc="ST" qchg="true" />
		<DA name="t" bType="Timestamp" fc="ST" />
		<DA name="ctlModel" type="CtlModels" bType="Enum" fc="CF" />
	</DOType>
	<DOType id="INS_1_Beh" cdc="INS">
		<DA name="stVal" bType="INT32" fc="ST" dchg="true" />
		<DA name="q" bType="Quality" fc="ST" qchg="true" />
		<DA name="t" bType="Timestamp" fc="ST" />
	</DOType>
	<DOType id="LPL_1_NamPlt" cdc="LPL">
		<DA name="vendor" bType="VisString255" fc="DC" />
		<DA name="swRev" bType="VisString255" fc="DC" />
		<DA name="d" bType="VisString255" fc="DC" />
		<DA name="configRev" bType="VisString255" fc="DC" />
		<DA name="ldNs" bType="VisString255" fc="EX" />
	</DOType>
	<DOType id="DPL_1_PhyNam" cdc="DPL">
		<DA name="vendor" bType="VisString255" fc="DC" />
	</DOType>
	<DOType id="SPS_1_Proxy" cdc="SPS">
		<DA name="stVal" bType="BOOLEAN" fc="ST" dchg="true" />
		<DA name="q" bType="Quality" fc="ST" qchg="true" />
		<DA name="t" bType="Timestamp" fc="ST" />
	</DOType>
	<DOType id="LPL_2_NamPlt" cdc="LPL">
		<DA name="vendor" bType="VisString255" fc="DC" />
		<DA name="swRev" bType="VisString255" fc="DC" />
		<DA name="d" bType="VisString255" fc="DC" />
	</DOType>
	<DOType id="ENC_1_Mod" cdc="ENC">
		<DA name="q" bType="Quality" fc="ST" qchg="true" />
		<DA name="t" bType="Timestamp" fc="ST" />
		<DA name="ctlModel" type="CtlModels" bType="Enum" fc="CF" />
	</DOType>
	<DOType id="ENS_1_Beh" cdc="ENS">
		<DA name="stVal" bType="INT32" fc="ST" dchg="true" />
		<DA name="q" bType="Quality" fc="ST" qchg="true" />
		<DA name="t" bType="Timestamp" fc="ST" />
	</DOType>
	<DOType id="MV_1_TotWh" cdc="MV">
		<DA name="mag" type="AnalogueValue_1" bType="Struct" fc="MX" dchg="true" />
		<DA name="q" bType="Quality" fc="MX" qchg="true" />
		<DA name="t" bType="Timestamp" fc="MX" />
	</DOType>
	<DOType id="SCA_1_SchdAbsTm" cdc="SCA">
		<DA name="val" bType="FLOAT32" count="255" fc="SP" dchg="true" />
		<DA name="time" bType="Timestamp" count="255" fc="SP" dchg="true" />
	</DOType>
	<DOType id="ING_1_SchdId" cdc="ING" />
	<DAType id="AnalogueValue_1">
		<BDA name="f" bType="FLOAT32" />
	</DAType>
	<EnumType id="CtlModels">
		<EnumVal ord="0">status-only</EnumVal>
		<EnumVal ord="1">direct-with-normal-security</EnumVal>
		<EnumVal ord="2">sbo-with-normal-security</EnumVal>
		<EnumVal ord="3">direct-with-enhanced-security</EnumVal>
		<EnumVal ord="4">sbo-with-enhanced-security</EnumVal>
	</EnumType>
</DataTypeTemplates>
</SCL>
