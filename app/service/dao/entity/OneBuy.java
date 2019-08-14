package com.asvital.expand.entity;

import javax.persistence.*;
import java.io.Serializable;

/**
 * Created by sixf on 2016/7/3.
 */
@Entity
@Table
public class OneBuy implements Serializable {

    @Id
    @Column(updatable = false)
    @GeneratedValue(strategy =GenerationType.IDENTITY)
    private Long id;

    @Column(updatable = false)
    private String productID;
    private Long date;

    private Integer number;
    private Integer views;

    @OneToOne(cascade = CascadeType.REFRESH,fetch = FetchType.EAGER)
    @JoinColumn(name = "winerID")
    private Subscriber winer;

    private String code;

    private Integer status=Status.ORDER;//状态
    private Long overDate;//停止下单时间

    public Long getId() {
        return id;
    }
    public void setId(Long id) {
        this.id = id;
    }



    public Long getDate() {
        return date;
    }

    public void setDate(Long date) {
        this.date = date;
    }



    public String getCode() {
        return code;
    }

    public void setCode(String code) {
        this.code = code;
    }

    public Long getOverDate() {
        return overDate;
    }

    public void setOverDate(Long overDate) {
        this.overDate = overDate;
    }

    public Integer getStatus() {
        return status;
    }

    public void setStatus(Integer status) {
        this.status = status;
    }

    public String getProductID() {
        return productID;
    }

    public void setProductID(String productID) {
        this.productID = productID;
    }

    public Integer getNumber() {
        return number;
    }

    public void setNumber(Integer number) {
        this.number = number;
    }

    public Subscriber getWiner() {
        return winer;
    }

    public void setWiner(Subscriber winer) {
        this.winer = winer;
    }

    public Integer getViews() {
        return views;
    }

    public void setViews(Integer views) {
        this.views = views;
    }

    public static class Status{
        public static final Integer ORDER=0;
        public static final Integer OVER_ORDER=1;
        public static final Integer WAIT_OPEN=2;
        public static final Integer OPEN=3;
    }
}
