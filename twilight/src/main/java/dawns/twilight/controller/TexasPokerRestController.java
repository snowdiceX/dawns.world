package dawns.twilight.controller;

import dawns.twilight.common.base.BaseRestController;
import dawns.twilight.common.web.JsonResult;
import dawns.twilight.dao.model.TexasPoker;
import dawns.twilight.dao.model.TexasPokerExample;
import dawns.twilight.service.TexasPokerService;
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
@RequestMapping("/texaspoker")
@Api(value = "TexasPoker api", description = "TexasPoker api",tags = {"TexasPoker"})
public class TexasPokerRestController extends BaseRestController{
    @Autowired
    private TexasPokerService texasPokerService;

    @ApiOperation(value="新增TexasPoker",tags = {"TexasPoker"})
    @RequestMapping(value = "", method = RequestMethod.POST)
    public JsonResult<Integer> add(HttpServletRequest request, @RequestBody TexasPoker texasPoker) {
        return new JsonResult<>(texasPokerService.insert(texasPoker));
    }

    @ApiOperation(value="根据id更新TexasPoker",tags = {"TexasPoker"})
    @RequestMapping(value = "", method = RequestMethod.PUT)
    public JsonResult<Integer> update(HttpServletRequest request, @RequestBody TexasPoker texasPoker) {
        return new JsonResult<>(texasPokerService.updateByPrimaryKey(texasPoker));
    }

    @ApiOperation(value="根据id查询TexasPoker",tags = {"TexasPoker"})
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    public JsonResult<TexasPoker> get(HttpServletRequest request, @PathVariable("id") Integer id) {
        TexasPoker texasPoker=texasPokerService.selectByPrimaryKey(id);
        if(texasPoker!=null){
            return new JsonResult<>(texasPoker);
        }else{
            return new JsonResult<>(HttpStatus.NOT_FOUND);
        }
    }

    @ApiOperation(value="根据id删除TexasPoker",tags = {"TexasPoker"})
    @RequestMapping(value = "/{id}", method = RequestMethod.DELETE)
    public JsonResult<Integer> delete(HttpServletRequest request, @PathVariable("id") Integer id) {
        return new JsonResult<>(texasPokerService.deleteByPrimaryKey(id));
    }

    @ApiOperation(value="分页查询TexasPoker",tags = {"TexasPoker"})
    @RequestMapping(value = "", method = RequestMethod.GET)
    public JsonResult<List<TexasPoker>> page(HttpServletRequest request,
                                            @RequestParam(value = "pageNum") Integer pageNum,
                                            @RequestParam(value = "pageSize") Integer pageSize) {
        TexasPokerExample texasPokerExample=new TexasPokerExample();
        texasPokerExample.setOrderByClause("id");
        return new JsonResult<>(texasPokerService.selectByExampleForStartPage(texasPokerExample, pageNum,pageSize));
    }
}
