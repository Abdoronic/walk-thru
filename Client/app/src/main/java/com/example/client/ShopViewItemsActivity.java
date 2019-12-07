package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

public class ShopViewItemsActivity extends AppCompatActivity {

    private Button ordersButton;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_shop_view_items);

        ordersButton = findViewById(R.id.ordersButton);

        ordersButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(ShopViewItemsActivity.this, ShopViewOrdersActivity.class);
                startActivity(i);
            }
        });
    }
}
