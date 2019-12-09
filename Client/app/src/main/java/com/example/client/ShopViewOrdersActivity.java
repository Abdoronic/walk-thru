package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
import android.widget.Toast;

import com.android.volley.DefaultRetryPolicy;
import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonArrayRequest;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.Volley;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class ShopViewOrdersActivity extends AppCompatActivity {
    private static final String TAG = "ShopViewOrdersActivity";
    private Button itemsButton;
    private String[] orderData;
    private JSONObject[] orderDataJSON;
    private String[] deliveredData;
    private JSONObject[] deliveredDataJSON;
    private ArrayAdapter<String> myAdapter;
    private ArrayAdapter<String> deliveredAdapter;
    private ListView ordersListView;
    private ListView deliveredOrdersListView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_shop_view_orders);

        itemsButton = findViewById(R.id.itemsButton);
        ordersListView = findViewById(R.id.ordersListView);
        deliveredOrdersListView = findViewById(R.id.deliveredOrdersList);

        itemsButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(ShopViewOrdersActivity.this, ShopViewItemsActivity.class);
                i.putExtra("id",getIntent().getIntExtra("id",-1));
                i.putExtra("name",getIntent().getStringExtra("name"));
                i.putExtra("location",getIntent().getStringExtra("location"));
                i.putExtra("adminUsername",getIntent().getStringExtra("adminUsername"));
                i.putExtra("adminPassword",getIntent().getStringExtra("adminPassword"));
                startActivity(i);
            }
        });

        RequestQueue queue = Volley.newRequestQueue(ShopViewOrdersActivity.this);

        String url = "http://10.0.2.2:8000/shops/"+getIntent().getIntExtra("id",-1)+"/viewPendingOrders";
        JsonArrayRequest jsonArrayRequest = new JsonArrayRequest(Request.Method.GET, url, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                try {
                    orderData = new String[response.length()];
                    orderDataJSON = new JSONObject[response.length()];
                    for(int i=0;i<response.length();i++){
                        JSONObject orderJSON = response.getJSONObject(i);
                        String order = "Order Number: " + orderJSON.getString("id") + "\n"
                                + "Delivered: " + orderJSON.getString("delivered") + "\n"
                                +  "Price: " +  orderJSON.getString("price") + "\n"
                                + "Date: " + orderJSON.getString("date") + '\n';
                        orderData[i]=order;
                        orderDataJSON[i]=orderJSON;
                    }
                    myAdapter = new ArrayAdapter<String>(ShopViewOrdersActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, orderData);
                    ordersListView.setAdapter(myAdapter);
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                try {
                    JSONObject errData =new JSONObject(new String(error.networkResponse.data));
                    Toast.makeText(getApplicationContext(),errData.getString("error"),Toast.LENGTH_LONG).show();

                } catch (JSONException e) {
                    e.printStackTrace();
                }
                error.printStackTrace();
            }
        });
        jsonArrayRequest.setRetryPolicy(new DefaultRetryPolicy(
                20000,
                DefaultRetryPolicy.DEFAULT_MAX_RETRIES,
                DefaultRetryPolicy.DEFAULT_BACKOFF_MULT));
        queue.add(jsonArrayRequest);

        String url2 = "http://10.0.2.2:8000/shops/"+getIntent().getIntExtra("id",-1)+"/viewDeliveredOrders";
        JsonArrayRequest jsonArrayRequest2 = new JsonArrayRequest(Request.Method.GET, url2, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                try {
                    deliveredData = new String[response.length()];
                    deliveredDataJSON = new JSONObject[response.length()];
                    for(int i=0;i<response.length();i++){
                        JSONObject orderJSON = response.getJSONObject(i);
                        String order = "Order Number: " + orderJSON.getString("id") + "\n"
                                + "Delivered: " + orderJSON.getString("delivered") + "\n"
                                +  "Price: " +  orderJSON.getString("price") + "\n"
                                + "Date: " + orderJSON.getString("date") + '\n';
                        deliveredData[i]=order;
                        deliveredDataJSON[i]=orderJSON;
                    }
                    Log.d(TAG, "onResponse: "+deliveredData[0]);
                    deliveredAdapter = new ArrayAdapter<String>(ShopViewOrdersActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, deliveredData);
                    deliveredOrdersListView.setAdapter(deliveredAdapter);
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        }, new Response.ErrorListener() {
            @Override
            public void onErrorResponse(VolleyError error) {
                try {
                    JSONObject errData =new JSONObject(new String(error.networkResponse.data));
                    Toast.makeText(getApplicationContext(),errData.getString("error"),Toast.LENGTH_LONG).show();

                } catch (JSONException e) {
                    e.printStackTrace();
                }
                error.printStackTrace();
            }
        });
        jsonArrayRequest2.setRetryPolicy(new DefaultRetryPolicy(
                20000,
                DefaultRetryPolicy.DEFAULT_MAX_RETRIES,
                DefaultRetryPolicy.DEFAULT_BACKOFF_MULT));
        queue.add(jsonArrayRequest2);

        ordersListView.setOnItemClickListener(new AdapterView.OnItemClickListener() {
            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                try {
                    RequestQueue queue = Volley.newRequestQueue(ShopViewOrdersActivity.this);
                    String url = "http://10.0.2.2:8000/shops/"+getIntent().getIntExtra("id",-1)+"/deliverOrder/"+orderDataJSON[position].getInt("id");
                    JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.PUT, url, null, new Response.Listener<JSONObject>() {
                        @Override
                        public void onResponse(JSONObject response) {
                            Toast.makeText(getApplicationContext(), "Order Delivered!",Toast.LENGTH_LONG).show();
                            finish();
                            overridePendingTransition(0, 0);
                            startActivity(getIntent());
                            overridePendingTransition(0, 0);
                        }
                    }, new Response.ErrorListener() {
                        @Override
                        public void onErrorResponse(VolleyError error) {
                            Toast.makeText(getApplicationContext(), "Order Delivered!",Toast.LENGTH_LONG).show();
                            error.printStackTrace();
                            finish();
                            overridePendingTransition(0, 0);
                            startActivity(getIntent());
                            overridePendingTransition(0, 0);
                        }
                    });
                    jsonObjectRequest.setRetryPolicy(new DefaultRetryPolicy(
                            20000,
                            DefaultRetryPolicy.DEFAULT_MAX_RETRIES,
                            DefaultRetryPolicy.DEFAULT_BACKOFF_MULT));
                    queue.add(jsonObjectRequest);
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        });
    }
}
