package dawns.twilight.common.web;

import lombok.Data;

@Data
public class RegisterTransaction {

	private String chain;
	private String Token;
	private String Contract;
	private String From;
	private String To;
	private String Amount;
	private String Txhash;
	private String Height;
}
