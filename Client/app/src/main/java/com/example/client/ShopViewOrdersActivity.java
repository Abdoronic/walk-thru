package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;

public class ShopViewOrdersActivity extends AppCompatActivity {

    private Button itemsButton;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_shop_view_orders);

        itemsButton = findViewById(R.id.itemsButton);

        itemsButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(ShopViewOrdersActivity.this, ShopViewItemsActivity.class);
                startActivity(i);
            }
        });
    }
}
