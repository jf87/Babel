package edu.berkeley.babel;

import android.os.Bundle;
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

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.Iterator;

import edu.berkeley.babel.util.JSONArrayHttpGetTask;
import edu.berkeley.babel.util.JSONArrayHttpGetTask.onJSONArrayHttpGetRespondedListener;
import edu.berkeley.babel.util.JSONObjectHttpPostTask;
import edu.berkeley.babel.util.JSONObjectHttpPostTask.onJSONObjectHttpPostRespondedListener;

public class MainActivity extends ActionBarActivity {

    private TextView mTypeText;
    private Spinner mTypeSpinner;
    private ArrayAdapter<String> mTypeSpinnerAdapter;
    private ListView mAttributeList;
    private AttributeListAdapter mAttributeListAdapter;
    private Button mStartButton;
    private TextView mActionText;
    private TextView mActionDesc;

    private boolean mBusy = false;
    private JSONArray mMetadataArray = null;
    private JSONObject mCurMetadata = null;

    /**
     * response to the AsyncTask that GETs metadata from server
     */
    private class GetMetadataArrayListener implements onJSONArrayHttpGetRespondedListener {
        @Override
        public void onJSONArrayHttpGetResponded(JSONArray response) {
            setUIEnabled(true);
            mBusy = false;

            if (response == null) {
                // TODO show error message
                return;
            }

            mMetadataArray = response;
            refreshType();
        }
    }

    /**
     * response to the AsyncTask that POSTs user-updated metadata to server
     */
    private class PostMetadataListener implements onJSONObjectHttpPostRespondedListener {
        @Override
        public void onJSONObjectHttpPostResponded(JSONObject response) {
            setUIEnabled(true);
            mBusy = false;

            startInstruction();
        }
    }

    /**
     * response to user selecting the type spinner
     */
    private class TypeSpinnerListener implements AdapterView.OnItemSelectedListener {
        @Override
        public void onItemSelected(AdapterView<?> parent, View view,
                                   int pos, long id) {
            // An item was selected. You can retrieve the selected item using
            // parent.getItemAtPosition(pos)
            if (mBusy) { // this should not happen
                return;
            }

            updateCurMetadataRef();
            refreshAttributes();
        }

        @Override
        public void onNothingSelected(AdapterView<?> parent) {
            // Another interface callback
        }
    }

    /**
     * response to user pressing the start button
     */
    private class StartOnClickListener implements View.OnClickListener {
        @Override
        public void onClick(View v) {
            // Perform action on click
            if (mBusy) { // this should not happen
                return;
            }

            updateCurMetadataFromUI();
            postCurMetadataToServer();
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

        mAttributeListAdapter.clear();

        Iterator<String> iter = mCurMetadata.keys();

        while (iter.hasNext()) {
            String name = iter.next();
            if (name.equals("kind") || name.equals("sequence")) {
                continue;
            }
            try {
                String value = mCurMetadata.getString(name);
                mAttributeListAdapter.add(new AttributeListAdapter.Pair<>(name, value));
            } catch (JSONException e) {
                e.printStackTrace();
            }
        }

        mAttributeListAdapter.notifyDataSetChanged();
    }

    /**
     * Update the mCurMetadata based on the type spinner selection
     */
    private void updateCurMetadataRef() {
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
        JSONArrayHttpGetTask httpGetTask = new JSONArrayHttpGetTask(new GetMetadataArrayListener());
        URL url = getHttpURL(getString(R.string.server), Integer.parseInt(getString(R.string.port)), getString(R.string.types_path));
        mBusy = true;
        setUIEnabled(false);
        httpGetTask.execute(url);
    }

    /**
     * Update mCurMetadata from UI
     */
    private void updateCurMetadataFromUI() {
        if (mCurMetadata == null) {
            return;
        }

        try {
            int count = mAttributeListAdapter.getCount();
            for (int i = 0; i < count; i++) {
                AttributeListAdapter.Pair<String, String> attr = mAttributeListAdapter.getItem(i);
                mCurMetadata.put(attr.first, attr.second);
            }
        } catch (JSONException e) {
            e.printStackTrace();
        }
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
     * start showing instruction to user to control the device
     */
    private void startInstruction() {

    }

    /**
     * Enable/disable all UI components
     */
    private void setUIEnabled(boolean enabled) {
        mTypeText.setEnabled(enabled);
        mTypeSpinner.setEnabled(enabled);
        mAttributeListAdapter.setEnabled(enabled);
        mStartButton.setEnabled(enabled);
        mActionText.setEnabled(enabled);
        mActionDesc.setEnabled(enabled);
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
        mAttributeListAdapter = new AttributeListAdapter(this);
        mAttributeList.setAdapter(mAttributeListAdapter);

        mStartButton = (Button) findViewById(R.id.start_button);
        mStartButton.setOnClickListener(new StartOnClickListener());

        mActionText = (TextView) findViewById(R.id.action_text);
        mActionDesc = (TextView) findViewById(R.id.action_desc);

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
