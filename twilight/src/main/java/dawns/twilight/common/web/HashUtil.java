package dawns.twilight.common.web;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.UnsupportedEncodingException;
import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 *
 */
public class HashUtil {

    public final static String md5(String content) {
        //用于加密的字符
        char[] md5String = {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
                'A', 'B', 'C', 'D', 'E', 'F'};
        try {
            //使用平台的默认字符集将此 String 编码为 byte序列，并将结果存储到一个新的 byte数组中
            byte[] btInput = content.getBytes();

            //信息摘要是安全的单向哈希函数，它接收任意大小的数据，并输出固定长度的哈希值。
            MessageDigest mdInst = MessageDigest.getInstance("MD5");

            //MessageDigest对象通过使用 update方法处理数据， 使用指定的byte数组更新摘要
            mdInst.update(btInput);

            // 摘要更新之后，通过调用digest（）执行哈希计算，获得密文
            byte[] md = mdInst.digest();

            // 把密文转换成十六进制的字符串形式
            int j = md.length;
            char[] str = new char[j * 2];
            int k = 0;
            for (int i = 0; i < j; i++) {   //  i = 0
                byte byte0 = md[i];  //95
                str[k++] = md5String[byte0 >>> 4 & 0xf];    //    5
                str[k++] = md5String[byte0 & 0xf];   //   F
            }

            //返回经过加密后的字符串
            return new String(str);

        } catch (Exception e) {
            return null;
        }
    }
    

    /**
     * 获取文件的md5值 ，有可能不是32位
     * @param filePath	文件路径
     * @return
     * @throws FileNotFoundException
     */
    public static String md5HashCode(String filePath) throws FileNotFoundException{  
        FileInputStream fis = new FileInputStream(filePath);  
        try {  
        	//拿到一个MD5转换器,如果想使用SHA-1或SHA-256，则传入SHA-1,SHA-256  
            MessageDigest md = MessageDigest.getInstance("MD5"); 
            
            //分多次将一个文件读入，对于大型文件而言，比较推荐这种方式，占用内存比较少。
            byte[] buffer = new byte[1024];  
            int length = -1;  
            while ((length = fis.read(buffer, 0, 1024)) != -1) {  
                md.update(buffer, 0, length);  
            }  
            fis.close();
            //转换并返回包含16个元素字节数组,返回数值范围为-128到127
  			byte[] md5Bytes  = md.digest();
            BigInteger bigInt = new BigInteger(1, md5Bytes);//1代表绝对值 
            return bigInt.toString(16);//转换为16进制
        } catch (Exception e) {  
            e.printStackTrace();  
            return "";  
        }  
    }
    
    public static String sha256(String str){
    	 MessageDigest messageDigest;
    	 String encodeStr = "";
    	 try {
    	  messageDigest = MessageDigest.getInstance("SHA-256");
    	  messageDigest.update(str.getBytes("UTF-8"));
    	  encodeStr = byte2Hex(messageDigest.digest());
    	 } catch (NoSuchAlgorithmException e) {
    	  e.printStackTrace();
    	 } catch (UnsupportedEncodingException e) {
    	  e.printStackTrace();
    	 }
    	 return encodeStr;
    }

    /**
    * 将byte转为16进制
    * @param bytes
    * @return
    */
    private static String byte2Hex(byte[] bytes){
     StringBuffer stringBuffer = new StringBuffer();
     String temp = null;
     for (int i=0;i<bytes.length;i++){
      temp = Integer.toHexString(bytes[i] & 0xFF);
      if (temp.length()==1){
      //1得到一位的进行补0操作
      stringBuffer.append("0");
      }
      stringBuffer.append(temp);
     }
     return stringBuffer.toString();
    }
}
