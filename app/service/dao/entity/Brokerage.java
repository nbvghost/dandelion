package com.asvital.expand.entity;

import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by sixf on 3/13/2016.
 */
@Entity
@Table
public class Brokerage implements Serializable {

    @Id
    @Column()
    @GeneratedValue(generator = "paymentableGenerator")
    @GenericGenerator(name = "paymentableGenerator", strategy = "uuid")
    private String id;
    @Column(updatable = false)
    private Long date;
    @Column(updatable = false)
    private String payShopID;
    @Column(updatable = false)
    private String receiveShopID;
    @Column(updatable = false)
    private Integer ratio;
    @Column(updatable = false)
    private Integer amount;
    @Column(updatable = false)
    private String orderID;
    @Column(columnDefinition = "Boolean default false")
    private Boolean balance = false;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public Long getDate() {
        return date;
    }

    public void setDate(Long date) {
        this.date = date;
    }


    public Integer getRatio() {
        return ratio;
    }

    public void setRatio(Integer ratio) {
        this.ratio = ratio;
    }

    public Integer getAmount() {
        return amount;
    }

    public void setAmount(Integer amount) {
        this.amount = amount;
    }

    public String getOrderID() {
        return orderID;
    }

    public void setOrderID(String orderID) {
        this.orderID = orderID;
    }

    public Boolean getBalance() {
        return balance;
    }

    public void setBalance(Boolean balance) {
        this.balance = balance;
    }

    public String getPayShopID() {
        return payShopID;
    }

    public void setPayShopID(String payShopID) {
        this.payShopID = payShopID;
    }

    public String getReceiveShopID() {
        return receiveShopID;
    }

    public void setReceiveShopID(String receiveShopID) {
        this.receiveShopID = receiveShopID;
    }
}
