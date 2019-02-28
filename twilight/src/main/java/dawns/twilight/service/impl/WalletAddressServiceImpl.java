package dawns.twilight.service.impl;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import dawns.twilight.common.annotation.BaseService;
import dawns.twilight.common.base.BaseServiceImpl;
import dawns.twilight.dao.mapper.WalletAddressMapper;
import dawns.twilight.dao.model.WalletAddress;
import dawns.twilight.dao.model.WalletAddressExample;
import dawns.twilight.service.WalletAddressService;

/**
* WalletAddressService实现
* Auto generate on 2019/2/28.
*/
@Service
@Transactional
@BaseService
public class WalletAddressServiceImpl extends BaseServiceImpl<WalletAddressMapper, WalletAddress, WalletAddressExample> implements WalletAddressService {
    @Autowired
    WalletAddressMapper walletAddressMapper;

    @Override
    protected WalletAddressMapper getMapper(){
        return walletAddressMapper;
    }
}
