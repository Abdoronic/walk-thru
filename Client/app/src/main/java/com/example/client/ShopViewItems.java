package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.os.Bundle;
import android.widget.TextView;

public class ShopViewItems extends AppCompatActivity {

    private TextView welcomeShopTextView;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_shop_view_items);

        welcomeShopTextView = findViewById(R.id.welcomeShopTextView);
        welcomeShopTextView.setText("Welcome, "+ getIntent().getStringExtra("name"));
    }
}
