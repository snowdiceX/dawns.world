package dawns.twilight.common.web;

import lombok.Data;

@Data
public class RegisterTransaction {
	private String chain;
	private String token;
	private String contract;
	private String from;
	private String to;
	private String amount;
	private String txhash;
}
