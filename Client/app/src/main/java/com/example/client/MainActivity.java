package com.example.client;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.VolleyError;
import com.android.volley.toolbox.JsonObjectRequest;
import com.android.volley.toolbox.Volley;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.util.HashMap;
import java.util.Map;

public class MainActivity extends AppCompatActivity {
    private Button signInButton;
    private TextView emailTextView;
    private TextView passwordTextView;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        Button signUp = (Button) findViewById(R.id.signUpBtn);
        signInButton = findViewById(R.id.signInBtn);
        emailTextView = findViewById(R.id.emailText);
        passwordTextView = findViewById((R.id.passwordText));

        signUp.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent startIntent = new Intent(getApplicationContext(),SignUpActivity.class);
                startActivity(startIntent);
            }
        });

        Button Admin = (Button) findViewById(R.id.adminLogIn);
        Admin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent startIntent = new Intent(getApplicationContext(),AdminActivity.class);
                startActivity(startIntent);
            }
        });

        signInButton.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                RequestQueue queue = Volley.newRequestQueue(MainActivity.this);

                Map<String, String> params = new HashMap<String, String>();
                params.put("email", emailTextView.getText().toString());
                params.put("password", passwordTextView.getText().toString());

                String url = "http://10.0.2.2:8000/customers/login";
                JsonObjectRequest jsonObjectRequest = new JsonObjectRequest(Request.Method.POST, url, new JSONObject(params), new Response.Listener<JSONObject>() {
                    @Override
                    public void onResponse(JSONObject response) {
                        Intent i = new Intent(getApplicationContext(),CustomerViewShops.class);
                        try {
                            i.putExtra("firstName",response.getString("firstName"));
                            i.putExtra("id",response.getInt("id"));

                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        startActivity(i);
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
                queue.add(jsonObjectRequest);
            }
        });
    }
}
