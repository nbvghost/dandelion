package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by sixf on 2016/9/14.
 */
@Entity
@Table
public class OneBuyProduct implements Serializable {
    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;
    @OneToOne(cascade = CascadeType.REFRESH,fetch = FetchType.EAGER)
    @JoinColumn(name = "productID")
    private Products products;

    @OneToOne(cascade = CascadeType.REFRESH,fetch = FetchType.EAGER)
    @JoinColumn(name = "oneBuyID")
    private OneBuy oneBuy;

    @Column(updatable = false)
    private String shopID;

    private Integer unit=Unit.YiYuan;//1元商品
    @Column(updatable = false)
    private Long date;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public Products getProducts() {
        return products;
    }

    public void setProducts(Products products) {
        this.products = products;
    }

    public OneBuy getOneBuy() {
        return oneBuy;
    }

    public void setOneBuy(OneBuy oneBuy) {
        this.oneBuy = oneBuy;
    }

    public String getShopID() {
        return shopID;
    }

    public void setShopID(String shopID) {
        this.shopID = shopID;
    }

    public Integer getUnit() {
        return unit;
    }

    public void setUnit(Integer unit) {
        this.unit = unit;
    }

    public Long getDate() {
        return date;
    }

    public void setDate(Long date) {
        this.date = date;
    }


    public static class Unit{
        public static  final Integer YiYuan = 100;//一元
        public static  final Integer ShiYuan = 1000;//十元
        public static  final Integer BaiYuan = 10000;//百元
    }

}
