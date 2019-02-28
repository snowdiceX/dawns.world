package dawns.twilight.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import dawns.twilight.common.annotation.BaseService;
import dawns.twilight.common.base.BaseServiceImpl;
import dawns.twilight.dao.mapper.LoginUserMapper;
import dawns.twilight.dao.model.LoginUser;
import dawns.twilight.dao.model.LoginUserExample;
import dawns.twilight.service.LoginUserService;

/**
* LoginUserService实现
* Auto generate on 2019/2/28.
*/
@Service
@Transactional
@BaseService
public class LoginUserServiceImpl extends BaseServiceImpl<LoginUserMapper, LoginUser, LoginUserExample> implements LoginUserService {
    @Autowired
    LoginUserMapper loginUserMapper;

    @Override
    protected LoginUserMapper getMapper(){
        return loginUserMapper;
    }
}
