package dawns.twilight.common.web;

import lombok.Data;

@Data
public class ResponseWallet {

	private String network;
	
	private String token;
	
	private String address;
	
	private String txid;
	
	private int height;
}
