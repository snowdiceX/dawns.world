package dawns.twilight;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.ApplicationContext;
import org.springframework.scheduling.annotation.EnableScheduling;
import org.springframework.transaction.annotation.EnableTransactionManagement;

import dawns.twilight.common.web.ApplicationContextUtil;
import springfox.documentation.swagger2.annotations.EnableSwagger2;

@SpringBootApplication
@EnableConfigurationProperties
@EnableTransactionManagement
@EnableSwagger2
@EnableScheduling
public class App {
	
    public static void main(String[] args) {
    	ApplicationContext ctx = SpringApplication.run(App.class, args);
		ApplicationContextUtil.setContext(ctx);
    }
}
