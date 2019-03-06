package dawns.twilight.common.cazimi;

import org.springframework.stereotype.Service;

import com.sun.jna.Library;
import com.sun.jna.Native;

import lombok.extern.slf4j.Slf4j;

@SuppressWarnings("all")
@Service
@Slf4j
public class CazimiService {

	// hyperledger fabric api
	public interface Cazimi extends Library {

		public void initJNI();
		public String chaincodeInvoke(String chainID,
				String chaincodeID, String args);
	}
	
	private static String libpath = System.getProperty("user.dir")
			+ System.getProperty("file.separator") + "lib" + System.getProperty("file.separator")
			+ "cazimi.so";
	static String ERR_NULLJNA="{\"code\": \"404\", \"height\": 0, \"reason\": \"jna library is null\"}";
	private static Cazimi api ;
	static {
		api = (Cazimi) Native.loadLibrary(libpath, Cazimi.class);
		api.initJNI();
	}
	
	// 注册代理（托管）钱包地址
	public String ChaincodeInvoke(String chainID, String chaincodeID, String args) {
		log.info("ChaincodeInvoke: "+chainID+"; "+chaincodeID+"; "+args);
		String ret = ERR_NULLJNA;
		if (api != null) {
			ret = api.chaincodeInvoke(chainID, chaincodeID, args);
		}
		log.info(ret);
		return 	ret;
	}
	
	public String registerWallet(String userId, String token, String network){
		return ChaincodeInvoke("dawns.world", "wallet",
				"{\"Func\":\"register\", \"Args\":[\""+userId+"\", \""+token+"\", \""+network+"\"]}");
	}
	
	public static void main(String[] args) {
		String ret = null;
		CazimiService cazimi = new CazimiService();
		ret = cazimi.ChaincodeInvoke("dawns.world", "wallet", "{\"Func\":\"register\"}");
		System.out.println("cazimi: "+ret);
	}
}
