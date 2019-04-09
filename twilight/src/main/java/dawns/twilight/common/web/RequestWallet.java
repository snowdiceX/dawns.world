package dawns.twilight.common.web;

import lombok.Data;

@Data
public class RequestWallet {

	private String network;
	
	private String token;
	
	private String pass;
	
	private String address;
	
	private String txid;
	
	private String height;
}
