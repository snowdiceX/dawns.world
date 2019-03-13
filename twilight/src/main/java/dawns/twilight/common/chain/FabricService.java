package dawns.twilight.common.chain;

import org.springframework.stereotype.Service;

import com.sun.jna.Library;
import com.sun.jna.Native;

import lombok.extern.slf4j.Slf4j;

@SuppressWarnings("all")
@Service
@Slf4j
public class FabricService {

	// hyperledger fabric api
	public interface FabrciApi extends Library {

		public String chaincodeInvoke(String channelID,
				String chaincodeID, String args);
		public String chaincodeQuery(String channelID,
				String chaincodeID, String args);
	}
	
	private static String libpath = System.getProperty("user.dir")
			+ System.getProperty("file.separator") + "lib" + System.getProperty("file.separator")
			+ "fabric.so";
	static String ERR_NULLJNA="{\"code\": \"404\", \"height\": 0, \"reason\": \"jna library is null\"}";
	private static FabrciApi api ;
	static {
		api = (FabrciApi) Native.loadLibrary(libpath, FabrciApi.class);
	}
	
	public String ChaincodeInvoke(String chainID, String chaincodeID, String args) {
		log.info("ChaincodeInvoke: "+chainID+"; "+chaincodeID+"; "+args);
		String ret = ERR_NULLJNA;
		if (api != null) {
			ret = api.chaincodeInvoke(chainID, chaincodeID, args);
		}
		log.info(ret);
		return 	ret;
	}
	
	public String ChaincodeQuery(String chainID, String chaincodeID, String args) {
		log.info("ChaincodeQuery: "+chainID+"; "+chaincodeID+"; "+args);
		String ret = ERR_NULLJNA;
		if (api != null) {
			ret = api.chaincodeQuery(chainID, chaincodeID, args);
		}
		log.info(ret);
		return 	ret;
	}
	
	// 注册代理（托管）钱包地址
	public String registerWallet(String userId, String token, String network){
		return ChaincodeInvoke("orgchannel", "wallet",
				"{\"Func\":\"register\", \"Args\":[\""+userId+"\", \""+token+"\", \""+network+"\"]}");
	}
	
	public static void main(String[] args) {
		String ret = null;
		FabricService api = new FabricService();
//		ret = api.ChaincodeInvoke("dawns.world", "wallet", "{\"Func\":\"register\"}");
//		System.out.println("cazimi: "+ret);
		ret = api.ChaincodeInvoke("orgchannel", "wallet", "{\"Func\":\"query\"}");
		System.out.println("fabric: "+ret);
	}
}
