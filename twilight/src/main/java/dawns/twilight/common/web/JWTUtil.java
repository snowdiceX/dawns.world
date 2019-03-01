package dawns.twilight.common.web;

import com.auth0.jwt.JWT;
import com.auth0.jwt.JWTVerifier;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.exceptions.JWTDecodeException;
import com.auth0.jwt.interfaces.DecodedJWT;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;

import java.util.Date;

@Configuration
public class JWTUtil {

    private static long EXPIRE_TIME = 31 * 24 * 60 * 60 * 1000;

    @Value("${token.expire.time}")
    public void setKafkaZk(long expireTime) {
        EXPIRE_TIME = expireTime;
    }

    /**
     * 校验token是否正确
     *
     * @param token  密钥
     * @param secret 用户的密码
     * @return 是否正确
     */
    public static boolean verify(String token, Integer uid, String secret) {
        try {
            Algorithm algorithm = Algorithm.HMAC256(secret);
            JWTVerifier verifier = JWT.require(algorithm)
                                      .withClaim("uid", uid)
                                      .build();
//            DecodedJWT jwt = verifier.verify(token);
            verifier.verify(token);
            return true;
        } catch (Exception exception) {
            return false;
        }

    }

    /**
     * 获得token中的信息无需secret解密也能获得
     *
     * @return token中包含的用户名
     */
    public static Integer getUserId(String token) {
        try {
            DecodedJWT jwt = JWT.decode(token);
            return jwt.getClaim("uid").asInt();
        } catch (JWTDecodeException e) {
            return null;
        }
    }

    /**
     * @param uid    用户id
     * @param secret 用户的密码
     * @return 加密的token
     */
    public static String sign(Integer uid, String secret) {
        Date date = new Date(System.currentTimeMillis() + EXPIRE_TIME);
        Algorithm algorithm = Algorithm.HMAC256(secret);
        return JWT.create()
                  .withClaim("uid", uid)
                  .withExpiresAt(date)
                  .sign(algorithm);
    }
}
