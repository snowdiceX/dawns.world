package dawns.twilight.common.web;

import lombok.Data;

@Data
public class FundsRequest {

	private String fundsKey;
	
	private String walletAddress;
	
	private String chain;
	
	private String token;
	
	private String amount;
}
