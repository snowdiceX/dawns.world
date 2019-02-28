package dawns.twilight.common.web;

import java.io.Serializable;
import java.util.Date;

import org.springframework.http.HttpStatus;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import lombok.Data;

/**
 * @author gaoxiang
 * @Description: JSON 基本类
 */
@JsonIgnoreProperties(ignoreUnknown = true)
@Data
public class JsonResult<T> implements Serializable {

    private static final long serialVersionUID = 1L;

    private Integer code = HttpStatus.OK.value();

    private String message = HttpStatus.OK.getReasonPhrase();

    private Date time = new Date();

    private T data;

    public JsonResult(T data) {
        this.data = data;
    }

    public JsonResult(HttpStatus httpStatus) {
        this.code = httpStatus.value();
        this.message = httpStatus.getReasonPhrase();
    }

    public JsonResult(HttpStatus httpStatus, T data) {
        this.code = httpStatus.value();
        this.message = httpStatus.getReasonPhrase();
        this.data = data;
    }

    public JsonResult(Integer code,String message) {
        this.code = code;
        this.message = message;
    }
}
