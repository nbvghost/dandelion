package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by sixf on 2016/7/25.
 */
@Entity
@Table
public class Express implements Serializable {
    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;

    private String s_name;
    private String s_tel;
    private String s_region;
    private String s_address;

    private String r_name;
    private String r_tel;
    private String r_region;
    private String r_address;

    private String des;//描述
    private String photos;//图片
    private String remark;//备注

    private String code;//快递编号

    private String shopID;//所在的商铺ID

    @ManyToOne(cascade = CascadeType.REFRESH,fetch = FetchType.EAGER)
    @JoinColumn(name = "executorID")
    private Subscriber executor;//处理人，执行者ID
    private Long allocationDate;//分配处理人时间


    private String provider="yuantong";//默认圆通快递，后期可能还有顺风等其它服务商
    @Column(updatable = false)
    private String userID;

    private Boolean selfVisit;//是否自己上门寄件

    @Column(updatable = false)
    private Long createDate;

    @OneToOne(cascade = CascadeType.ALL,fetch = FetchType.EAGER)
    @JoinColumn(name = "ordersID")
    private Orders orders;

    public Express() {
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getS_name() {
        return s_name;
    }

    public void setS_name(String s_name) {
        this.s_name = s_name;
    }

    public String getS_tel() {
        return s_tel;
    }

    public void setS_tel(String s_tel) {
        this.s_tel = s_tel;
    }

    public String getS_region() {
        return s_region;
    }

    public void setS_region(String s_region) {
        this.s_region = s_region;
    }

    public String getS_address() {
        return s_address;
    }

    public void setS_address(String s_address) {
        this.s_address = s_address;
    }

    public String getR_name() {
        return r_name;
    }

    public void setR_name(String r_name) {
        this.r_name = r_name;
    }

    public String getR_tel() {
        return r_tel;
    }

    public void setR_tel(String r_tel) {
        this.r_tel = r_tel;
    }

    public String getR_region() {
        return r_region;
    }

    public void setR_region(String r_region) {
        this.r_region = r_region;
    }

    public String getR_address() {
        return r_address;
    }

    public void setR_address(String r_address) {
        this.r_address = r_address;
    }

    public String getDes() {
        return des;
    }

    public void setDes(String des) {
        this.des = des;
    }

    public String getPhotos() {
        return photos;
    }

    public void setPhotos(String photos) {
        this.photos = photos;
    }

    public String getRemark() {
        return remark;
    }

    public void setRemark(String remark) {
        this.remark = remark;
    }

    public Boolean getSelfVisit() {
        return selfVisit;
    }

    public void setSelfVisit(Boolean selfVisit) {
        this.selfVisit = selfVisit;
    }

    public Long getCreateDate() {
        return createDate;
    }

    public void setCreateDate(Long createDate) {
        this.createDate = createDate;
    }

    public String getUserID() {
        return userID;
    }

    public void setUserID(String userID) {
        this.userID = userID;
    }

    public String getProvider() {
        return provider;
    }

    public void setProvider(String provider) {
        this.provider = provider;
    }

    public String getShopID() {
        return shopID;
    }

    public void setShopID(String shopID) {
        this.shopID = shopID;
    }



    public Orders getOrders() {
        return orders;
    }

    public void setOrders(Orders orders) {
        this.orders = orders;
    }

    public Subscriber getExecutor() {
        return executor;
    }

    public void setExecutor(Subscriber executor) {
        this.executor = executor;
    }

    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public Long getAllocationDate() {
        return allocationDate;
    }

    public void setAllocationDate(Long allocationDate) {
        this.allocationDate = allocationDate;
    }

    public static class Provider{
        public static final String yuantong = "yuantong";
        public static String GetChineseName(String provider){
            switch (provider){
                case yuantong:
                    return "圆通";
                default:
                    return "圆通";

            }
        }
    }
}
