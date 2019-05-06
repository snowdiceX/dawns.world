package dawns.twilight.common.web;

import lombok.Data;

@Data
public class TxRegister {
	private String chain;
	private String token;
	private String address;
	private String amount;
	private String gasUsed;
	private String gasPrice;
	private TxInfo info;
}
