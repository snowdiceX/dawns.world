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

		public String chaincodeInvoke(String channelID, String chaincodeID, String args);
		public String chaincodeQuery(String channelID, String chaincodeID, String args);
		public String registerToken(String chain, String token);
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
	public String RegisterWallet(String pass, String chain, String token){
		return this.ChaincodeInvoke("orgchannel", "wallet",
				new StringBuilder("{\"Func\":\"create\", \"Args\":[\"wallet\",\"")
				.append(chain).append("\",\"")
				.append(token).append("\",\"")
				.append(pass).append("\"]}").toString());
	}
	
	public String ImportWallet(String address, String chain, String token, String height) {
		return this.ChaincodeInvoke("orgchannel", "wallet",
				new StringBuilder("{\"Func\":\"register\", \"Args\":[\"wallet\",\"")
				.append(address).append("\",\"")
				.append(chain).append("\",\"")
				.append(token).append("\",\"")
				.append(height).append("\"]}").toString());
	}
	
	public String RegisterToken(String chain, String token) {
		String ret = ERR_NULLJNA;
		if (api != null) {
			ret =  api.registerToken(chain, token);
		}
		log.info(ret);
		return 	ret;
	}
	
	public String RegisterFunds(String baseChain, String baseToken, String chain, String token) {
		String ret = ERR_NULLJNA;
		if (api != null) {
			StringBuilder buf = new StringBuilder();
			buf.append("{\"Func\":\"register\", \"Args\":[\"funds\", \"")
			.append(baseChain).append("\", \"")
			.append(baseToken).append("\", \"")
			.append(chain).append("\", \"")
			.append(token).append("\"]}");
			ret =  api.chaincodeInvoke("orgchannel", "wallet", buf.toString());
		}
		log.info(ret);
		return 	ret;
	}
	
	public String FundsDeposit(String fundsKey,
			String walletAddress, String chain,
			String token, String amount) {
		String ret = ERR_NULLJNA;
		if (api != null) {
			StringBuilder buf = new StringBuilder();
			buf.append("{\"Func\":\"funds\", \"Args\":[\"deposit\", \"")
			.append(fundsKey).append("\", \"")
			.append(walletAddress).append("\", \"")
			.append(chain).append("\", \"")
			.append(token).append("\", \"")
			.append(amount).append("\"]}");
			ret =  api.chaincodeInvoke("orgchannel", "wallet", buf.toString());
		}
		log.info(ret);
		return 	ret;
	}
	
	public String FundsWithdraw(String fundsKey,
			String walletAddress, String chain,
			String token, String amount) {
		String ret = ERR_NULLJNA;
		if (api != null) {
			StringBuilder buf = new StringBuilder();
			buf.append("{\"Func\":\"funds\", \"Args\":[\"withdraw\", \"")
			.append(fundsKey).append("\", \"")
			.append(walletAddress).append("\", \"")
			.append(chain).append("\", \"")
			.append(token).append("\", \"")
			.append(amount).append("\"]}");
			ret =  api.chaincodeInvoke("orgchannel", "wallet", buf.toString());
		}
		log.info(ret);
		return 	ret;
	}
	
	public static void main(String[] args) {
		String ret = null;
		FabricService api = new FabricService();
//		ret = api.ChaincodeInvoke("orgchannel", "wallet", "{\"Func\":\"register\"}");
//		System.out.println("cazimi: "+ret);
		ret = api.ChaincodeInvoke("orgchannel", "wallet", "{\"Func\":\"query\"}");
		System.out.println("fabric: "+ret);
	}
}
