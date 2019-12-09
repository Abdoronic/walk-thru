package com.example.client;

public class Card {
    private String title;
    private String desc;
    private String price;

    public Card() {
        this.title = "";
        this.desc = "";
        this.price = "";
    }

    public Card(String title, String desc, String price) {
        this.title = title;
        this.desc = desc;
        this.price = price;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getPrice() {
        return price;
    }

    public void setPrice(String price) {
        this.price = price;
    }

}
