package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;

public class CustomerViewOrdersActivity extends AppCompatActivity {

    private Button shopButton;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_orders);

        shopButton = findViewById(R.id.shopsButton);
        if(getIntent().getStringExtra("activity").equals("items")){
            shopButton.setText("Items");
            shopButton.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    Intent i = new Intent(CustomerViewOrdersActivity.this, CustomerViewItemsActivity.class);
                    startActivity(i);
                }
            });
        }
        else{
            shopButton.setText("Shops");
            shopButton.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    Intent i = new Intent(CustomerViewOrdersActivity.this, CustomerViewShopsActivity.class);
                    startActivity(i);
                }
            });
        }


    }
}
