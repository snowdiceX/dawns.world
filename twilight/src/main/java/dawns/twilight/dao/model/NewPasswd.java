package dawns.twilight.dao.model;

import lombok.Data;

@Data
public class NewPasswd {

    /**
     * 
     */
    private String email;

    /**
     * 
     */
    private String authCode;
    
    /**
     * 
     */
    private String newPassword;
}
