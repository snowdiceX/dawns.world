package dawns.twilight.common.web;

import org.springframework.context.ApplicationContext;

public class ApplicationContextUtil {
	private static ApplicationContext context;  
	  
 
	public static void setContext(ApplicationContext context) {
		ApplicationContextUtil.context = context;
	}


	public static ApplicationContext getApplicationContext() {  
        return context;  
    }  
}
