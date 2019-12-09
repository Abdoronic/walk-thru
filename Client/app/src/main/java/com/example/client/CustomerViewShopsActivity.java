package com.example.client;

import androidx.appcompat.app.ActionBar;
import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.graphics.Color;
import android.graphics.Typeface;
import android.os.Bundle;
import android.view.Gravity;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.RequestBody;


public class CustomerViewShopsActivity extends AppCompatActivity {


    private Button shopsButton, ordersButton;

    private String[] data;
    private JSONObject[] dataJSON;
    private JSONObject[] shopsJSON;
    private ArrayAdapter<String> shopsListAdapter, ordersListAdapter;
    private ListView listView;
    private boolean shopsLoaded = false, ordersLoaded = false, ordersDisplayed = false;

    public void loadShops() {
        OkHttpClient client = new OkHttpClient();

        String url = getString(R.string.BASE_URL) + "/shops";

        okhttp3.Request request = new okhttp3.Request.Builder()
                .url(url)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                e.printStackTrace();
            }

            @Override
            public void onResponse(Call call, okhttp3.Response response) throws IOException {
                if (response.isSuccessful()) {
                    try {
                        JSONArray res = new JSONArray(response.body().string());
                        data = new String[res.length()];
                        shopsJSON = new JSONObject[res.length()];
                        for (int i = 0; i < res.length(); i++) {
                            JSONObject shopJSON = res.getJSONObject(i);
                            String shop = "Shop Name: " + shopJSON.getString("name") + "\n"
                                    + "Shop Location: " + shopJSON.getString("location");
                            data[i] = shop;
                            shopsJSON[i] = shopJSON;
                        }
                        shopsListAdapter = new ArrayAdapter<String>(CustomerViewShopsActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, data);
                        shopsLoaded = true;
                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                } else {
                    try {
                        final JSONObject error = new JSONObject(response.body().string());
                        CustomerViewShopsActivity.this.runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                try {
                                    Toast.makeText(getApplicationContext(),
                                            error.getString("error"),
                                            Toast.LENGTH_SHORT).show();
                                } catch (Exception e) {
                                    e.printStackTrace();
                                }
                            }
                        });
                        shopsLoaded = true;
                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                }
            }
        });
    }

    public void loadOrders() {
        OkHttpClient client = new OkHttpClient();

        String url = getString(R.string.BASE_URL) + "/customers/" + getIntent().getIntExtra("id", -1) + "/viewOrders";

        okhttp3.Request request = new okhttp3.Request.Builder()
                .url(url)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                e.printStackTrace();
            }

            @Override
            public void onResponse(Call call, okhttp3.Response response) throws IOException {
                if (response.isSuccessful()) {
                    try {
                        JSONArray res = new JSONArray(response.body().string());
                        data = new String[res.length()];
                        dataJSON = new JSONObject[res.length()];
                        for (int i = 0; i < res.length(); i++) {
                            JSONObject orderJSON = res.getJSONObject(i);
                            String order = "Order Number: " + orderJSON.getString("id") + "\n"
                                    + "Delivered: " + orderJSON.getString("delivered") + "\n"
                                +  "Price: " +  orderJSON.getString("price") + "\n"
                                    + "Date: " + orderJSON.getString("date") + '\n';
                            data[i] = order;
                            dataJSON[i] = orderJSON;
                        }
                        ordersListAdapter = new ArrayAdapter<String>(CustomerViewShopsActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, data);
                        ordersLoaded = true;
                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                } else {
                    try {
                        final JSONObject error = new JSONObject(response.body().string());
                        CustomerViewShopsActivity.this.runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                try {
                                    Toast.makeText(getApplicationContext(),
                                            error.getString("error"),
                                            Toast.LENGTH_SHORT).show();
                                } catch (Exception e) {
                                    e.printStackTrace();
                                }
                            }
                        });
                        ordersLoaded = true;
                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                }
            }
        });
    }

    public void viewShop(final int position) {
        if(ordersDisplayed) return;
        OkHttpClient client = new OkHttpClient();

        String url = getString(R.string.BASE_URL) + "/customers/" +
                getIntent().getIntExtra("id", -1) + "/createOrder";

        String bodyJson = new StringBuilder()
                .append("{")
                .append("\"delivered\":").append("false")
                .append("}")
                .toString();

        RequestBody body = RequestBody.create(
                MediaType.parse("application/json; charset=utf-8"),
                bodyJson
        );
        okhttp3.Request request = new okhttp3.Request.Builder()
                .url(url)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                e.printStackTrace();
            }

            @Override
            public void onResponse(Call call, okhttp3.Response response) throws IOException {
                if (response.isSuccessful()) {
                    try {
                        final JSONObject data = new JSONObject(response.body().string());
                        CustomerViewShopsActivity.this.runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                try {
                                    int shopID = shopsJSON[position].getInt("id");
                                    Intent i = new Intent(getApplicationContext(), CustomerViewItemsActivity.class);

                                    i.putExtra("id", getIntent().getIntExtra("id",-1));
                                    i.putExtra("firstName", getIntent().getStringExtra("firstName"));
                                    i.putExtra("lastName", getIntent().getStringExtra("lastName"));
                                    i.putExtra("email", getIntent().getStringExtra("email"));
                                    i.putExtra("password", getIntent().getStringExtra("password"));
                                    i.putExtra("creditCardNumber", getIntent().getStringExtra("creditCardNumber"));
                                    i.putExtra("creditCardExpiryDate", getIntent().getStringExtra("creditCardExpiryDate"));
                                    i.putExtra("creditCardCVV", getIntent().getIntExtra("creditCardCVV",-1));

                                    i.putExtra("shopID", shopID);
                                    i.putExtra("orderID", data.getInt("id"));

                                    startActivity(i);
                                } catch (JSONException e) {
                                    e.printStackTrace();
                                }
                            }
                        });
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                } else {
                    try {
                        final JSONObject error = new JSONObject(response.body().string());
                        CustomerViewShopsActivity.this.runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                try {
                                    Toast.makeText(getApplicationContext(), error.getString("error"), Toast.LENGTH_SHORT).show();
                                } catch (Exception e) {
                                    e.printStackTrace();
                                }
                            }
                        });
                    } catch (JSONException e) {
                        e.printStackTrace();
                    }
                }
            }
        });
    }

    public void showShops() {
        while(!shopsLoaded);
        listView.setAdapter(shopsListAdapter);
        listView.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, final int position, long id) {
                viewShop(position);
            }
        });

    }

    public void showSOrders() {
        while(!ordersLoaded);
        listView.setAdapter(ordersListAdapter);
    }

    public void setTitle(String title) {
        getSupportActionBar().setHomeButtonEnabled(true);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);
        TextView textView = new TextView(this);
        textView.setText(title);
        textView.setTextSize(30);
        textView.setTypeface(getResources().getFont(R.font.pacifico), Typeface.NORMAL);
        textView.setLayoutParams(new LinearLayout.LayoutParams(LinearLayout.LayoutParams.FILL_PARENT,
                LinearLayout.LayoutParams.WRAP_CONTENT));
        textView.setGravity(Gravity.CENTER);
        textView.setTextColor(getResources().getColor(R.color.textColorPrimary));
        getSupportActionBar().setDisplayOptions(ActionBar.DISPLAY_SHOW_CUSTOM);
        getSupportActionBar().setCustomView(textView);
        getSupportActionBar().show();
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_shops);

        setTitle("Walk Thru");

        shopsButton = findViewById(R.id.shopsButton);
        ordersButton = findViewById(R.id.ordersButton);
        listView = findViewById(R.id.listView);

        ordersButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                shopsButton.setBackground(getDrawable(R.drawable.greybutton));
                ordersButton.setBackground(getDrawable(R.drawable.orangebutton));
                ordersDisplayed = true;
                showSOrders();
            }
        });
        shopsButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                shopsButton.setBackground(getDrawable(R.drawable.orangebutton));
                ordersButton.setBackground(getDrawable(R.drawable.greybutton));
                ordersDisplayed = false;
                showShops();
            }
        });
        loadShops();
        loadOrders();
        showShops();
    }
}
