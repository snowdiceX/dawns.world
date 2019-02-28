package dawns.twilight.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import dawns.twilight.common.annotation.BaseService;
import dawns.twilight.common.base.BaseServiceImpl;
import dawns.twilight.dao.mapper.PublicRandomMapper;
import dawns.twilight.dao.model.PublicRandom;
import dawns.twilight.dao.model.PublicRandomExample;
import dawns.twilight.service.PublicRandomService;

/**
* PublicRandomService实现
* Auto generate on 2019/2/28.
*/
@Service
@Transactional
@BaseService
public class PublicRandomServiceImpl extends BaseServiceImpl<PublicRandomMapper, PublicRandom, PublicRandomExample> implements PublicRandomService {
    @Autowired
    PublicRandomMapper publicRandomMapper;

    @Override
    protected PublicRandomMapper getMapper(){
        return publicRandomMapper;
    }
}
