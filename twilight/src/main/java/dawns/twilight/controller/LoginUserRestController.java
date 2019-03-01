package dawns.twilight.controller;

import java.util.List;

import javax.servlet.http.HttpServletRequest;

import org.apache.shiro.authz.annotation.RequiresAuthentication;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.HashUtil;
import dawns.twilight.common.web.JWTUtil;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.common.web.LoginToken;
import dawns.twilight.dao.model.LoginUser;
import dawns.twilight.dao.model.LoginUserExample;
import dawns.twilight.dao.model.NewPasswd;
import dawns.twilight.dao.model.Signin;
import dawns.twilight.service.LoginUserService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;

@Controller
@RequestMapping("")
@Api(value = "Login api", tags = {"Login"})
public class LoginUserRestController extends BaseRestController{
	
    @Autowired
    private LoginUserService loginUserService;

    @ApiOperation(value="登录", tags = {"Login"})
    @RequestMapping(value = "/signin", method = RequestMethod.POST)
    public JsonResult<LoginToken> signin(HttpServletRequest request, @RequestBody Signin signin) {
    	LoginUserExample example = new LoginUserExample();
        example.createCriteria().andEmailEqualTo(signin.getEmail());
        LoginUser loginUser = loginUserService.selectFirstByExample(example);
        if(loginUser!=null && HashUtil.md5(signin.getPassword()).equals(loginUser.getPassword())){
        	String t = JWTUtil.sign(loginUser.getId(), loginUser.getPassword());
        	desensitize(loginUser);
        	LoginToken token = new LoginToken();
        	token.setToken(t);
        	token.setUser(loginUser);
            return new JsonResult<>(token);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }
    
    @ApiOperation(value="注册", tags = {"Login"})
    @PostMapping(value = "/signup")
    public JsonResult<LoginUser> signup(HttpServletRequest request, @RequestBody LoginUser loginUser) {
    	if (loginUser!=null) {
    		loginUser.setPassword(HashUtil.md5(loginUser.getPassword()));
    		int ret = loginUserService.insertSelective(loginUser);
    		if(ret>0){
    			desensitize(loginUser);
    			return new JsonResult<>(loginUser);
    		}
    	}

        return new JsonResult<>(HttpStatus.INTERNAL_SERVER_ERROR);
    }
    
    @ApiOperation(value="修改密码", tags = {"Login"})
    @PutMapping("/passwd/reset")
    public JsonResult<Boolean> update(HttpServletRequest request, @RequestBody NewPasswd passwd) {
        return new JsonResult<>(true);
    }

    @ApiOperation(value="分页查询（后台）",tags = {"Login"})
    @GetMapping(value = "/user/list")
    @RequiresAuthentication
    public JsonResult<List<LoginUser>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        LoginUserExample example = new LoginUserExample();
        example.setOrderByClause("id");
        List<LoginUser> list = loginUserService
        		.selectByExampleForStartPage(example, pageNum, pageSize);
        for (LoginUser u: list) {
        	desensitize(u);
        }
        return new JsonResult<>(list);
    }
    
    public void desensitize(LoginUser user) {
    	user.setPassword(null);
    }
}
