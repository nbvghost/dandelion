package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by SIX4 on 2014-11-06 .
 */
@Entity
@Table
public class PreItem implements Serializable {
    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;
    private String productID;
    private Long date;
    private String begin_timestamp="8:0";//秒杀
    private String end_timestamp="23:0";//秒杀
    private Integer total=15;//活动天数
    private Integer stock;
    private Float discount;
    private String targetID;
    private String type;
    private boolean disable = false;//目前转盘专用 ，为true 时，为转盘的参与奖

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getProductID() {
        return productID;
    }

    public void setProductID(String productID) {
        this.productID = productID;
    }


    public String getTargetID() {
        return targetID;
    }

    public void setTargetID(String targetID) {
        this.targetID = targetID;
    }

    public Integer getStock() {
        return stock;
    }

    public void setStock(Integer stock) {
        this.stock = stock;
    }

    public Float getDiscount() {
        return discount;
    }

    public void setDiscount(Float discount) {
        this.discount = discount;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public boolean isDisable() {
        return disable;
    }

    public void setDisable(boolean disable) {
        this.disable = disable;
    }

    public Long getDate() {
        return date;
    }

    public void setDate(Long date) {
        this.date = date;
    }

    public String getBegin_timestamp() {
        return begin_timestamp;
    }

    public void setBegin_timestamp(String begin_timestamp) {
        this.begin_timestamp = begin_timestamp;
    }

    public String getEnd_timestamp() {
        return end_timestamp;
    }

    public void setEnd_timestamp(String end_timestamp) {
        this.end_timestamp = end_timestamp;
    }

    public Integer getTotal() {
        return total;
    }

    public void setTotal(Integer total) {
        this.total = total;
    }

    public static class Type {

        public static final String preferential = "preferential";
        public static final String lottery = "lottery";
        public static final String seckill = "seckill";

    }
}
