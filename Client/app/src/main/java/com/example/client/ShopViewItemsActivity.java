package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonArrayRequest;
import com.android.volley.toolbox.Volley;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class ShopViewItemsActivity extends AppCompatActivity {

    private Button ordersButton;
    private String[] itemData;
    private JSONObject[] itemDataJSON;
    private ArrayAdapter<String> myAdapter;
    private ListView itemsListView;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_shop_view_items);

        ordersButton = findViewById(R.id.ordersButton);
        itemsListView = findViewById(R.id.itemsListView);

        ordersButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(ShopViewItemsActivity.this, ShopViewOrdersActivity.class);
                i.putExtra("id",getIntent().getIntExtra("id",-1));
                startActivity(i);
            }
        });

        RequestQueue queue = Volley.newRequestQueue(ShopViewItemsActivity.this);

        String url = "http://10.0.2.2:8000/shops/"+getIntent().getIntExtra("id",-1)+"/viewOfferedItems";
        JsonArrayRequest jsonArrayRequest = new JsonArrayRequest(Request.Method.GET, url, null, new Response.Listener<JSONArray>() {
            @Override
            public void onResponse(JSONArray response) {
                try {
                    itemData = new String[response.length()];
                    itemDataJSON = new JSONObject[response.length()];
                    for(int i=0;i<response.length();i++){
                        JSONObject itemJSON = response.getJSONObject(i);
                        String item = "Item Name: "+itemJSON.getString("name")+"\n"+"Item Type: "+itemJSON.getString("type")+"\n"+"Item Description: "+itemJSON.getString("description")+"\n"+"Item Price: "+itemJSON.getDouble("price")+"\n";
                        itemData[i]=item;
                        itemDataJSON[i]=itemJSON;
                    }
                    myAdapter = new ArrayAdapter<String>(ShopViewItemsActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, itemData);
                    itemsListView.setAdapter(myAdapter);
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
        queue.add(jsonArrayRequest);
    }
}
