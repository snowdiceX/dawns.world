package dawns.twilight.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import dawns.twilight.common.annotation.BaseService;
import dawns.twilight.common.base.BaseServiceImpl;
import dawns.twilight.dao.mapper.GamingMapper;
import dawns.twilight.dao.model.Gaming;
import dawns.twilight.dao.model.GamingExample;
import dawns.twilight.service.GamingService;

/**
* GamingService实现
* Auto generate on 2019/2/28.
*/
@Service
@Transactional
@BaseService
public class GamingServiceImpl extends BaseServiceImpl<GamingMapper, Gaming, GamingExample> implements GamingService {
    @Autowired
    GamingMapper gamingMapper;

    @Override
    protected GamingMapper getMapper(){
        return gamingMapper;
    }
}
