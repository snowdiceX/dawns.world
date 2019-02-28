package dawns.twilight.controller;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.Gaming;
import dawns.twilight.dao.model.GamingExample;
import dawns.twilight.service.GamingService;
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
@RequestMapping("/gaming")
@Api(value = "Gaming api", description = "Gaming api",tags = {"Gaming"})
public class GamingRestController extends BaseRestController{
    @Autowired
    private GamingService gamingService;

    @ApiOperation(value="新增Gaming",tags = {"Gaming"})
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody Gaming gaming) {
        return new JsonResult<>(gamingService.insert(gaming));
    }

    @ApiOperation(value="根据id更新Gaming",tags = {"Gaming"})
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody Gaming gaming) {
        return new JsonResult<>(gamingService.updateByPrimaryKey(gaming));
    }

    @ApiOperation(value="根据id查询Gaming",tags = {"Gaming"})
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<Gaming> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        Gaming gaming=gamingService.selectByPrimaryKey(id);
        if(gaming!=null){
            return new JsonResult<>(gaming);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除Gaming",tags = {"Gaming"})
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(gamingService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询Gaming",tags = {"Gaming"})
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<Gaming>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        GamingExample gamingExample=new GamingExample();
        gamingExample.setOrderByClause("id");
        return new JsonResult<>(gamingService.selectByExampleForStartPage(gamingExample, pageNum,pageSize));
    }
}
