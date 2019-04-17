package dawns.twilight.common.web;

import lombok.Data;

@Data
public class RegisterBlock {
	private String height;
	private RegisterTransaction[] transactions;
}