package dawns.twilight.common.web;

import dawns.twilight.dao.model.LoginUser;
import lombok.Data;

@Data
public class LoginToken {

	private String token;
	private LoginUser user;
}
