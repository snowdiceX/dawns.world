package dawns.twilight.common.web;

import lombok.Data;

@Data
public class ResponseToken {

	private String network;
	
	private String address;
	
	private String name;
	
	private String symbol;
	
	private int decimals;
		
	private String txid;
		
	private String height;
}
