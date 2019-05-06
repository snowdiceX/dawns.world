package dawns.twilight.common.web;

import lombok.Data;

@Data
public class TxInfo {

	private String contract;
	private String from;
	private String to;
	private String txhash;
}
