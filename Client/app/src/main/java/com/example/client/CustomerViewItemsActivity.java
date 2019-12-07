package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
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

public class CustomerViewItemsActivity extends AppCompatActivity {

    private Button orderButton;

    private String[] itemData;
    private JSONObject[] itemDataJSON;
    private ArrayAdapter<String> myAdapter;
    private ListView itemsListView;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_customer_view_items);

        orderButton = findViewById(R.id.ordersButton);
        itemsListView = findViewById(R.id.itemsListView);

        orderButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(CustomerViewItemsActivity.this, CustomerViewOrdersActivity.class);
                i.putExtra("activity","items");
                i.putExtra("shopID",getIntent().getIntExtra("shopID",-1));
                i.putExtra("orderID",getIntent().getIntExtra("orderID",-1));
                i.putExtra("customerID",getIntent().getIntExtra("customerID",-1));
                startActivity(i);
            }
        });




        RequestQueue queue = Volley.newRequestQueue(CustomerViewItemsActivity.this);

        String url = "http://10.0.2.2:8000/customers/viewItems/"+getIntent().getIntExtra("shopID",-1);
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
                    myAdapter = new ArrayAdapter<String>(CustomerViewItemsActivity.this, android.R.layout.simple_list_item_1, android.R.id.text1, itemData);
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



        itemsListView.setOnItemClickListener(new AdapterView.OnItemClickListener() {

            @Override
            public void onItemClick(AdapterView<?> parent, View view, int position, long id) {
                try {
                    int itemID = itemDataJSON[position].getInt("id");
                    int shopID = itemDataJSON[position].getInt("shopID");
                    Intent i = new Intent(getApplicationContext(),CustomerItemClickActivity.class);
                    i.putExtra("itemID",itemID);
                    i.putExtra("shopID",shopID);
                    i.putExtra("customerID",getIntent().getIntExtra("customerID",-1));
                    i.putExtra("orderID",getIntent().getIntExtra("orderID",-1));
                    startActivity(i);
                } catch (JSONException e) {
                    e.printStackTrace();
                }
            }
        });


    }
}
