package dawns.twilight.common.web;

import lombok.Data;

@Data
public class FundsRequest {

	private String fundsTokenKey;
	
	private String walletAddress;
	
	private String amount;
}
