package dawns.twilight.controller;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.LoginUser;
import dawns.twilight.dao.model.LoginUserExample;
import dawns.twilight.service.LoginUserService;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;

import javax.servlet.http.HttpServletRequest;
import java.util.List;

@Controller
@RequestMapping("/login")
@Api(value = "LoginUser api", description = "LoginUser api",tags = {"LoginUser"})
public class LoginUserRestController extends BaseRestController{
    @Autowired
    private LoginUserService loginUserService;

    @ApiOperation(value="新增LoginUser",tags = {"LoginUser"})
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody LoginUser loginUser) {
        return new JsonResult<>(loginUserService.insert(loginUser));
    }

    @ApiOperation(value="根据id更新LoginUser",tags = {"LoginUser"})
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody LoginUser loginUser) {
        return new JsonResult<>(loginUserService.updateByPrimaryKey(loginUser));
    }

    @ApiOperation(value="根据id查询LoginUser",tags = {"LoginUser"})
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<LoginUser> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        LoginUser loginUser=loginUserService.selectByPrimaryKey(id);
        if(loginUser!=null){
            return new JsonResult<>(loginUser);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除LoginUser",tags = {"LoginUser"})
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(loginUserService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询LoginUser",tags = {"LoginUser"})
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<LoginUser>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        LoginUserExample loginUserExample=new LoginUserExample();
        loginUserExample.setOrderByClause("id");
        return new JsonResult<>(loginUserService.selectByExampleForStartPage(loginUserExample, pageNum,pageSize));
    }
}
