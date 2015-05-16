package edu.berkeley.babel;

import android.os.Bundle;
import android.support.v7.app.ActionBarActivity;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.Spinner;
import android.widget.TextView;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;

import edu.berkeley.babel.util.JSONArrayHttpGetTask;
import edu.berkeley.babel.util.JSONArrayHttpGetTask.onJSONArrayHttpGetRespondedListener;


public class MainActivity extends ActionBarActivity {

    private Spinner mTypeSpinner;
    private ArrayAdapter<String> mTypeSpinnerAdapter;
    private TextView mActionView;
    private Button mStartButton;

    private boolean mBusy = false;
    private JSONArray mMetadata = null;

    private class MetadataListener implements onJSONArrayHttpGetRespondedListener {
        @Override
        public void onJSONArrayHttpGetResponded(JSONArray response) {
            if (response == null) {
                mBusy = false;
                return;
            }

            mMetadata = response;
            refreshType();
            mBusy = false;
        }
    }

    private class TypeSpinnerListener implements AdapterView.OnItemSelectedListener {
        @Override
        public void onItemSelected(AdapterView<?> parent, View view,
                                   int pos, long id) {
            // An item was selected. You can retrieve the selected item using
            // parent.getItemAtPosition(pos)
            if (mBusy) {
                return;
            }

            refreshAttributes();
        }

        @Override
        public void onNothingSelected(AdapterView<?> parent) {
            // Another interface callback
        }
    }

    private class StartOnClickListener implements View.OnClickListener {
        @Override
        public void onClick(View v) {
            // Perform action on click
            // TODO
            if (mBusy) {
                return;
            }

        }
    }

    private void refreshType() {
        mTypeSpinnerAdapter.clear();
        for (int i = 0; i < mMetadata.length(); i++) {
            try {
                JSONObject typeObj = mMetadata.getJSONObject(i);
                String typeName = typeObj.getString("kind");
                mTypeSpinnerAdapter.add(typeName);
            } catch (JSONException e) {
                e.printStackTrace();
            }
        }
        mTypeSpinnerAdapter.notifyDataSetChanged();
    }

    private void refreshAttributes() {
        // TODO pop up attributes
    }

    private void getMetadata() {
        JSONArrayHttpGetTask httpGetTask = new JSONArrayHttpGetTask(new MetadataListener());
        URL url = getHttpURL("130.226.142.195", 4444, "/api/types");
        mBusy = true;
        httpGetTask.execute(url);
    }

    private URL getHttpURL(String host, int port, String path) {
        URL url = null;
        try {
            URI uri = new URI("http", null, host, port, path, null, null);
            url = uri.toURL();
        } catch (URISyntaxException e) {
            e.printStackTrace();
        } catch (MalformedURLException e) {
            e.printStackTrace();
        }

        return url;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        // Set up UI
        mTypeSpinner = (Spinner) findViewById(R.id.type_spinner);
        mTypeSpinner.setOnItemSelectedListener(new TypeSpinnerListener());
        mTypeSpinnerAdapter = new ArrayAdapter<String>
                (this, R.layout.support_simple_spinner_dropdown_item);
        mTypeSpinnerAdapter.setDropDownViewResource(R.layout.support_simple_spinner_dropdown_item);
        mTypeSpinner.setAdapter(mTypeSpinnerAdapter);

        mActionView = (TextView) findViewById(R.id.action_text);

        mStartButton = (Button) findViewById(R.id.start_button);
        mStartButton.setOnClickListener(new StartOnClickListener());


        getMetadata();
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
