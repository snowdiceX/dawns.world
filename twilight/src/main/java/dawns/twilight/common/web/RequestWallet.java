package dawns.twilight.common.web;

import lombok.Data;

@Data
public class RequestWallet {

	private String network;
	
	private String token;
	
	private String name;
	
	private String key;
}
