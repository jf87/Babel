package edu.berkeley.babel;

import android.os.Bundle;
import android.os.Handler;
import android.support.v7.app.ActionBarActivity;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ListView;
import android.widget.Spinner;
import android.widget.TextView;
import android.widget.Toast;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.Iterator;

import edu.berkeley.babel.util.JSONObjectHttpGetTask;
import edu.berkeley.babel.util.JSONObjectHttpGetTask.onJSONObjectHttpGetRespondedListener;
import edu.berkeley.babel.util.JSONObjectHttpPostTask;
import edu.berkeley.babel.util.JSONObjectHttpPostTask.onJSONObjectHttpPostRespondedListener;
import edu.berkeley.babel.util.KeyValueListAdapter;

public class MainActivity extends ActionBarActivity {
    private TextView mTypeText;
    private Spinner mTypeSpinner;
    private ArrayAdapter<String> mTypeSpinnerAdapter;
    private ListView mAttributeList;
    private KeyValueListAdapter mKeyValueListAdapter;
    private Button mSendButton;
    private TextView mInfoText;
    private TextView mInfoDesc;

    private boolean mBusy = false;
    private JSONArray mMetadataArray = null;
    private JSONObject mCurMetadata = null;

    /**
     * response to the AsyncTask that GETs metadata from server
     */
    private class GetMetadataArrayListener implements onJSONObjectHttpGetRespondedListener {
        @Override
        public void onJSONObjectHttpGetResponded(JSONObject response) {
            if (response == null) {
                // user has to restart the app
                Toast.makeText(getApplicationContext(), getString(R.string.no_conn), Toast.LENGTH_LONG).show();
                return;
            }

            setUIEnabled(true);
            mBusy = false;

            try {
                mMetadataArray = response.getJSONArray("library");
            } catch (JSONException e) {
                e.printStackTrace();
            }
            refreshType();
        }
    }

    /**
     * response to the AsyncTask that POSTs user-updated metadata to server
     */
    private class PostMetadataListener implements onJSONObjectHttpPostRespondedListener {
        @Override
        public void onJSONObjectHttpPostResponded(JSONObject response) {
            if (response == null) {
                // only enable UI when connection fails
                setUIEnabled(true);
                mBusy = false;
                Toast.makeText(getApplicationContext(), getString(R.string.no_conn), Toast.LENGTH_LONG).show();
                return;
            }

            String resultStr = "";
            try {
                resultStr = response.getString("result");
                if (resultStr == null || resultStr.isEmpty()) {
                    Toast.makeText(getApplicationContext(), getString(R.string.server_error), Toast.LENGTH_LONG).show();
                }
            } catch (JSONException e) {
                e.printStackTrace();
            }

            mInfoDesc.setText(resultStr);
            setUIEnabled(true);
            mBusy = false;
        }
    }

    /**
     * response to user selecting the type spinner
     */
    private class TypeSpinnerListener implements AdapterView.OnItemSelectedListener {
        @Override
        public void onItemSelected(AdapterView<?> parent, View view, int pos, long id) {
            // An item was selected. You can retrieve the selected item using
            // parent.getItemAtPosition(pos)
            if (mBusy) { // this should not happen
                return;
            }

            updateCurMetadataRefOnType();
            refreshAttributes();
        }

        @Override
        public void onNothingSelected(AdapterView<?> parent) {
            // Another interface callback
        }
    }

    /**
     * response to user pressing the send button
     */
    private class SendOnClickListener implements View.OnClickListener {
        @Override
        public void onClick(View v) {
            // Perform action on click
            if (mBusy) { // this should not happen
                return;
            }

            postCurMetadataToServer();
        }
    }

    private class OnAttributeChangedListener implements KeyValueListAdapter.OnKeyValueChangedListener {
        @Override
        public void OnValueChanged(String key, String newValue) {
            // key won't be in mCurMetadata if these sequence happens:
            // 1. attribute is changed
            // 2. type is changed, then mCurMetadata is changed
            // 3. this OnValueChanged is called
            if (mCurMetadata == null || !mCurMetadata.has(key)) {
                return;
            }

            try {
                mCurMetadata.put(key, newValue);
            } catch (JSONException e) {
                e.printStackTrace();
            }
        }
    }

    /**
     * Refresh the types in type spinner using mMetadataArray
     */
    private void refreshType() {
        mTypeSpinnerAdapter.clear();
        for (int i = 0; i < mMetadataArray.length(); i++) {
            try {
                JSONObject typeObj = mMetadataArray.getJSONObject(i);
                String typeName = typeObj.getString("kind");
                mTypeSpinnerAdapter.add(typeName);
            } catch (JSONException e) {
                e.printStackTrace();
            }
        }
        mTypeSpinnerAdapter.notifyDataSetChanged();
    }

    /**
     * Refresh the attributes based on the current selected Type
     */
    private void refreshAttributes() {
        // dynamically populate UI
        if (mCurMetadata == null) {
            return;
        }

        mKeyValueListAdapter.clear();

        Iterator<String> iter = mCurMetadata.keys();

        while (iter.hasNext()) {
            String name = iter.next();
            if (name.equals("kind")) {
                continue;
            }
            try {
                String value = mCurMetadata.getString(name);
                mKeyValueListAdapter.add(new KeyValueListAdapter.Pair<>(name, value));
            } catch (JSONException e) {
                e.printStackTrace();
            }
        }

        mKeyValueListAdapter.notifyDataSetChanged();

        mInfoDesc.setText("");
    }

    /**
     * Update the mCurMetadata based on the type spinner selection
     */
    private void updateCurMetadataRefOnType() {
        String curType = mTypeSpinner.getSelectedItem().toString();

        // TODO optimize lookup by indexing by kind
        JSONObject metadata = null;
        for (int i = 0; i < mMetadataArray.length(); i++) {
            try {
                metadata = mMetadataArray.getJSONObject(i);
                String typeName = metadata.getString("kind");
                if (typeName.equals(curType)) {
                    break;
                }
            } catch (JSONException e) {
                e.printStackTrace();
            }
        }

        mCurMetadata = metadata;
    }

    /**
     * start an AsyncTask to GET the metadata array from the server
     */
    private void getMetadataArrayFromServer() {
        JSONObjectHttpGetTask httpGetTask = new JSONObjectHttpGetTask(new GetMetadataArrayListener());
        URL url = getHttpURL(getString(R.string.server), Integer.parseInt(getString(R.string.port)), getString(R.string.types_path));
        mBusy = true;
        setUIEnabled(false);
        httpGetTask.execute(url);
    }

    /**
     * start an AsyncTask to POST the user-updated metadata to the server
     */
    private void postCurMetadataToServer() {
        JSONObjectHttpPostTask httpPostTask = new JSONObjectHttpPostTask(new PostMetadataListener());
        URL url = getHttpURL(getString(R.string.server), Integer.parseInt(getString(R.string.port)), getString(R.string.link_path));

        mBusy = true;
        setUIEnabled(false);
        httpPostTask.execute(url, mCurMetadata);
    }

    /**
     * Enable/disable all UI components
     */
    private void setUIEnabled(boolean enabled) {
        mTypeText.setEnabled(enabled);
        mTypeSpinner.setEnabled(enabled);
        mKeyValueListAdapter.setEnabled(enabled);
        mSendButton.setEnabled(enabled);
        mInfoText.setEnabled(enabled);
        mInfoDesc.setEnabled(enabled);
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
        mTypeText = (TextView) findViewById(R.id.type_text);
        mTypeSpinner = (Spinner) findViewById(R.id.type_spinner);
        mTypeSpinner.setOnItemSelectedListener(new TypeSpinnerListener());
        mTypeSpinnerAdapter = new ArrayAdapter<String>
                (this, R.layout.support_simple_spinner_dropdown_item);
        mTypeSpinnerAdapter.setDropDownViewResource(R.layout.support_simple_spinner_dropdown_item);
        mTypeSpinner.setAdapter(mTypeSpinnerAdapter);

        mAttributeList = (ListView) findViewById(R.id.attributes_list);
        mKeyValueListAdapter = new KeyValueListAdapter(this);
        mKeyValueListAdapter.setOnKeyValueChangedListener(new OnAttributeChangedListener());
        mAttributeList.setAdapter(mKeyValueListAdapter);

        mSendButton = (Button) findViewById(R.id.send_button);
        mSendButton.setOnClickListener(new SendOnClickListener());

        mInfoText = (TextView) findViewById(R.id.info_text);
        mInfoDesc = (TextView) findViewById(R.id.info_desc);

        // get metadata from server to populate the type spinner
        getMetadataArrayFromServer();
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
