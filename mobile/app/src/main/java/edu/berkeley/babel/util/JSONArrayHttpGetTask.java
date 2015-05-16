/*
 * Copyright (c) 2013, Regents of the University of California
 * All rights reserved.
 * 
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions 
 * are met:
 * 
 *  - Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 *  - Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the
 *    distribution.
 * 
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS 
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL 
 * THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, 
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES 
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR 
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) 
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, 
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) 
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED 
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 */

/*
 * Author: Kaifei Chen <kaifei@eecs.berkeley.edu>
 */

package edu.berkeley.babel.util;

import java.io.BufferedInputStream;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;

import org.json.JSONArray;
import org.json.JSONException;

import android.os.AsyncTask;

/* 
 * HTTP Get Task posts with JSON Array and gets JSON Array returned
 */
public class JSONArrayHttpGetTask extends AsyncTask<Object, Void, JSONArray> {

    private final onJSONArrayHttpGetRespondedListener listener;

    public static interface onJSONArrayHttpGetRespondedListener {
        void onJSONArrayHttpGetResponded(JSONArray response);
    }

    public JSONArrayHttpGetTask(final onJSONArrayHttpGetRespondedListener listener) {
        this.listener = listener;
    }

    private InputStream httpGet(final HttpURLConnection connection,
                                final URL url) throws IOException {
        connection.setRequestMethod("GET");

        final InputStream in = new BufferedInputStream(
                connection.getInputStream());

        return in;
    }

    @Override
    protected JSONArray doInBackground(final Object... params) {
        final URL url = (URL) params[0];

        if (url == null) {
            return null;
        }

        // TODO reuse the connection
        HttpURLConnection connection = null;
        try {
            connection = (HttpURLConnection) url.openConnection();
            final InputStream in = httpGet(connection, url);

            final BufferedReader reader = new BufferedReader(
                    new InputStreamReader(in));

            String line;
            final StringBuilder sb = new StringBuilder();
            while ((line = reader.readLine()) != null) {
                sb.append(line + '\n');
            }
            reader.close();

            return new JSONArray(sb.toString());
        } catch (final IOException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        } catch (final JSONException e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        } finally {
            if (connection != null) {
                connection.disconnect();
            }
        }

        return null;
    }

    @Override
    protected void onPostExecute(final JSONArray response) {
        if (!isCancelled() && listener != null) {
            listener.onJSONArrayHttpGetResponded(response);
        }
    }
}