package com.example.client;

import androidx.appcompat.app.ActionBar;
import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.graphics.Typeface;
import android.os.Bundle;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.LinearLayout;
import android.widget.TextView;
import android.widget.Toast;

import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class MainActivity extends AppCompatActivity {
    private Button signInButton;
    private TextView emailTextView;
    private TextView passwordTextView;

    public void signIn(String email, String password) {
        OkHttpClient client = new OkHttpClient();

        String url = getString(R.string.BASE_URL) + "/customers/login";

        String bodyJson = new StringBuilder()
                .append("{")
                .append("\"email\":\"").append(email).append("\",")
                .append("\"password\":\"").append(password).append("\"}")
                .toString();

        RequestBody body = RequestBody.create(
                MediaType.parse("application/json; charset=utf-8"),
                bodyJson
        );
        Request request = new Request.Builder()
                .url(url)
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                e.printStackTrace();
            }
            @Override
            public void onResponse(Call call, Response response) throws IOException {
                if (response.isSuccessful()) {
                    try {
                        final JSONObject data = new JSONObject(response.body().string());
                        MainActivity.this.runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                Intent intent = new Intent(getApplicationContext(), CustomerViewShopsActivity.class);
                                try {
                                    intent.putExtra("firstName",data.getString("firstName"));
                                    intent.putExtra("id", data.getInt("id"));

                                } catch (JSONException e) {
                                    e.printStackTrace();
                                }
                                startActivity(intent);
                            }
                        });
                    } catch(Exception e) {
                        e.printStackTrace();
                    }
                } else {
                    try {
                        final JSONObject error = new JSONObject(response.body().string());
                        MainActivity.this.runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                try {
                                    Toast.makeText(getApplicationContext(), error.getString("error"),Toast.LENGTH_SHORT).show();
                                } catch(Exception e) {
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

    public void setTitle(String title){
        getSupportActionBar().setHomeButtonEnabled(true);
        getSupportActionBar().setDisplayHomeAsUpEnabled(true);
        TextView textView = new TextView(this);
        textView.setText(title);
        textView.setTextSize(30);
        textView.setTypeface(getResources().getFont(R.font.pacifico), Typeface.NORMAL);
        textView.setLayoutParams(new LinearLayout.LayoutParams(LinearLayout.LayoutParams.FILL_PARENT, LinearLayout.LayoutParams.WRAP_CONTENT));
        textView.setGravity(Gravity.CENTER);
        textView.setTextColor(getResources().getColor(R.color.textColorPrimary));
        getSupportActionBar().setDisplayOptions(ActionBar.DISPLAY_SHOW_CUSTOM);
        getSupportActionBar().setCustomView(textView);
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        getSupportActionBar().hide();

        Button signUp = findViewById(R.id.signUpBtn);
        signInButton = findViewById(R.id.customerSignUpButton);
        emailTextView = findViewById(R.id.emailText);
        passwordTextView = findViewById((R.id.passwordText));

        signUp.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent startIntent = new Intent(getApplicationContext(),SignUpActivity.class);
                startActivity(startIntent);
            }
        });

        Button admin = findViewById(R.id.adminLogIn);
        admin.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent startIntent = new Intent(getApplicationContext(),AdminActivity.class);
                startActivity(startIntent);
            }
        });

        signInButton.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                signIn(emailTextView.getText().toString(), passwordTextView.getText().toString());
            }
        });
    }
}
