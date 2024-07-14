import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface TelemetryQuery extends DataQuery {
  telemetry?: string;
  source?: string;
  withStreaming: boolean;
  graph: boolean;
}

export const defaultQuery: Partial<TelemetryQuery> = {
  telemetry: 'SpeedKmh',
  source: 'acc',
  withStreaming: true,
  graph: false,
};

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  playstationIP: string;
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}
