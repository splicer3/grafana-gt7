import { defaults } from 'lodash';

import React, { PureComponent, SyntheticEvent } from 'react';
import {InlineField, InlineSwitch, Select} from '@grafana/ui';
import {QueryEditorProps, SelectableValue} from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, TelemetryQuery } from './types';
import {dirtRallyOptions} from "./dirtRallyOptions";

export const sourceOptions = [
  { label: 'DiRT Rally 2.0', value: 'dirtRally2' },
  { label: 'Assetto Corsa Competizione', value: 'acc' },
];

type Props = QueryEditorProps<DataSource, TelemetryQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onTelemetryChange = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, telemetry: option.value });
    // executes the query
    onRunQuery();
  };

  onSourceChange = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, source: option.value });
    onRunQuery();
  };

  onWithStreamingChange = (event: SyntheticEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, withStreaming: event.currentTarget.checked });
    // executes the query
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { telemetry, source, withStreaming } = query;

    return (
      <div className="gf-form">
        <InlineField label="Source">
          <Select
              width={25}
              options={sourceOptions}
              value={source}
              onChange={this.onSourceChange}
              defaultValue={'dirtRally2'}
          />
        </InlineField>
        <Select
          width={25}
          options={dirtRallyOptions}
          value={telemetry}
          onChange={this.onTelemetryChange}
          defaultValue={'Time'}
        />
        <InlineField label="Enable streaming (v8+)">
          <InlineSwitch
              value={withStreaming || false}
              onChange={this.onWithStreamingChange}
              css=""
          />
        </InlineField>
      </div>
    );
  }
}
