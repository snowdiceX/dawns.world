package dawns.twilight.common.base;

import javax.servlet.http.HttpServletRequest;

import org.apache.shiro.ShiroException;
import org.springframework.beans.InvalidPropertyException;
import org.springframework.beans.TypeMismatchException;
import org.springframework.http.HttpStatus;
import org.springframework.http.converter.HttpMessageConversionException;
import org.springframework.validation.BindException;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.MissingServletRequestParameterException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import dawns.twilight.common.web.JsonResult;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@RestControllerAdvice
public abstract class BaseRestController {

    /**
     * 异常处理
     *
     * @return
     */
    @ResponseStatus(HttpStatus.METHOD_NOT_ALLOWED)
    @ExceptionHandler({ HttpRequestMethodNotSupportedException.class })
    public JsonResult handle405(HttpRequestMethodNotSupportedException e) {
        return new JsonResult(HttpStatus.METHOD_NOT_ALLOWED);
    }

    /**
     * 异常处理
     *
     * @return
     */
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    @ExceptionHandler({ TypeMismatchException.class, InvalidPropertyException.class, BindException.class,
                        MethodArgumentNotValidException.class, MissingServletRequestParameterException.class,
                        HttpMessageConversionException.class })
    public JsonResult handle400(Exception e) {
        return new JsonResult(HttpStatus.BAD_REQUEST.value(), e.getMessage());
    }

    /**
     * 异常处理
     *
     * @return
     */
    @ResponseStatus(HttpStatus.UNAUTHORIZED)
    @ExceptionHandler({ ShiroException.class })
    public JsonResult handle401(ShiroException e) {
        return new JsonResult(HttpStatus.UNAUTHORIZED.value(),"未登录或token已超时");
    }

    /**
     * 异常处理
     * @return
     */
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    @ExceptionHandler({ Exception.class })
    public JsonResult handException(HttpServletRequest request, Exception e) {
        log.warn(" ExceptionHandler!" + request.getRequestURI());
        log.error("", e);
        return new JsonResult(HttpStatus.INTERNAL_SERVER_ERROR);
    }
}
