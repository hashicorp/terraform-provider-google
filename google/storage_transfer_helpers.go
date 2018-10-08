package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/storagetransfer/v1"
	"time"
)

func expandDates(dates []interface{}) []*storagetransfer.Date {
	expandedDates := make([]*storagetransfer.Date, 0, len(dates))
	for _, raw := range dates {
		date := raw.([]interface{})
		expandedDates = append(expandedDates, &storagetransfer.Date{
			Day:   int64(extractFirstMapConfig(date)["day"].(int)),
			Month: int64(extractFirstMapConfig(date)["month"].(int)),
			Year:  int64(extractFirstMapConfig(date)["year"].(int)),
		})
	}
	return expandedDates
}

func flattenDates(dates []*storagetransfer.Date) []map[string]interface{} {
	datesSchema := make([]map[string]interface{}, 0, len(dates))
	for _, date := range dates {
		datesSchema = append(datesSchema, map[string]interface{}{
			"year":  date.Year,
			"month": date.Month,
			"day":   date.Day,
		})
	}
	return datesSchema
}

func expandTimeOfDays(times []interface{}) []*storagetransfer.TimeOfDay {
	expandedTimes := make([]*storagetransfer.TimeOfDay, 0, len(times))
	for _, raw := range times {
		time := raw.([]interface{})
		expandedTimes = append(expandedTimes, &storagetransfer.TimeOfDay{
			Hours:   int64(extractFirstMapConfig(time)["hours"].(int)),
			Minutes: int64(extractFirstMapConfig(time)["minutes"].(int)),
			Seconds: int64(extractFirstMapConfig(time)["seconds"].(int)),
			Nanos:   int64(extractFirstMapConfig(time)["nanos"].(int)),
		})
	}
	return expandedTimes
}

func flattenTimeOfDays(timeOfDays []*storagetransfer.TimeOfDay) []map[string]interface{} {
	timeOfDaysSchema := make([]map[string]interface{}, 0, len(timeOfDays))
	for _, timeOfDay := range timeOfDays {
		timeOfDaysSchema = append(timeOfDaysSchema, map[string]interface{}{
			"hours":   timeOfDay.Hours,
			"minutes": timeOfDay.Minutes,
			"seconds": timeOfDay.Seconds,
			"nanos":   timeOfDay.Nanos,
		})
	}
	return timeOfDaysSchema
}

func expandTransferSchedules(transferSchedules []interface{}) []*storagetransfer.Schedule {
	schedules := make([]*storagetransfer.Schedule, 0, len(transferSchedules))
	for _, raw := range transferSchedules {
		schedule := raw.(map[string]interface{})
		sched := &storagetransfer.Schedule{
			ScheduleStartDate: expandDates([]interface{}{schedule["schedule_start_date"]})[0],
		}

		if v, ok := schedule["schedule_end_date"]; ok && len(v.([]interface{})) > 0 {
			sched.ScheduleEndDate = expandDates([]interface{}{v})[0]
		}
		if v, ok := schedule["start_time_of_day"]; ok && len(v.([]interface{})) > 0 {
			sched.StartTimeOfDay = expandTimeOfDays([]interface{}{v})[0]
		}

		schedules = append(schedules, sched)
	}
	return schedules
}

func flattenTransferSchedules(transferSchedules []*storagetransfer.Schedule) []map[string][]map[string]interface{} {
	transferSchedulesSchema := make([]map[string][]map[string]interface{}, 0, len(transferSchedules))
	for _, transferSchedule := range transferSchedules {
		schedule := map[string][]map[string]interface{}{
			"schedule_start_date": flattenDates([]*storagetransfer.Date{transferSchedule.ScheduleStartDate}),
		}

		if transferSchedule.ScheduleEndDate != nil {
			schedule["schedule_end_date"] = flattenDates([]*storagetransfer.Date{transferSchedule.ScheduleEndDate})
		}

		if transferSchedule.StartTimeOfDay != nil {
			schedule["start_time_of_day"] = flattenTimeOfDays([]*storagetransfer.TimeOfDay{transferSchedule.StartTimeOfDay})
		}

		transferSchedulesSchema = append(transferSchedulesSchema, schedule)
	}
	return transferSchedulesSchema
}

func expandGcsData(gcsDatas []interface{}) []*storagetransfer.GcsData {
	datas := make([]*storagetransfer.GcsData, 0, len(gcsDatas))
	for _, raw := range gcsDatas {
		data := raw.(map[string]interface{})
		datas = append(datas, &storagetransfer.GcsData{
			BucketName: data["bucket_name"].(string),
		})
	}
	return datas
}

func flattenGcsData(gcsDatas []*storagetransfer.GcsData) []map[string]interface{} {
	datasSchema := make([]map[string]interface{}, 0, len(gcsDatas))
	for _, data := range gcsDatas {
		datasSchema = append(datasSchema, map[string]interface{}{
			"bucket_name": data.BucketName,
		})
	}
	return datasSchema
}

func expandAwsAccessKeys(awsAccessKeys []interface{}) []*storagetransfer.AwsAccessKey {
	datas := make([]*storagetransfer.AwsAccessKey, 0, len(awsAccessKeys))
	for _, raw := range awsAccessKeys {
		data := raw.(map[string]interface{})
		datas = append(datas, &storagetransfer.AwsAccessKey{
			AccessKeyId:     data["access_key_id"].(string),
			SecretAccessKey: data["secret_access_key"].(string),
		})
	}
	return datas
}

func flattenAwsAccessKeys(awsAccessKeys []*storagetransfer.AwsAccessKey) []map[string]interface{} {
	datasSchema := make([]map[string]interface{}, 0, len(awsAccessKeys))
	for _, data := range awsAccessKeys {
		datasSchema = append(datasSchema, map[string]interface{}{
			"access_key_id":     data.AccessKeyId,
			"secret_access_key": data.SecretAccessKey,
		})
	}
	return datasSchema
}

func expandAwsS3Data(awsS3Datas []interface{}) []*storagetransfer.AwsS3Data {
	datas := make([]*storagetransfer.AwsS3Data, 0, len(awsS3Datas))
	for _, raw := range awsS3Datas {
		data := raw.(map[string]interface{})
		datas = append(datas, &storagetransfer.AwsS3Data{
			BucketName:   data["bucket_name"].(string),
			AwsAccessKey: expandAwsAccessKeys(data["aws_access_key"].([]interface{}))[0],
		})
	}
	return datas
}

func flattenAwsS3Data(awsS3Datas []*storagetransfer.AwsS3Data) []map[string]interface{} {
	datasSchema := make([]map[string]interface{}, 0, len(awsS3Datas))
	for _, data := range awsS3Datas {
		datasSchema = append(datasSchema, map[string]interface{}{
			"bucket_name":    data.BucketName,
			"aws_access_key": data.AwsAccessKey,
		})
	}
	return datasSchema
}

func expandHttpData(httpDatas []interface{}) []*storagetransfer.HttpData {
	datas := make([]*storagetransfer.HttpData, 0, len(httpDatas))
	for _, raw := range httpDatas {
		data := raw.(map[string]interface{})
		datas = append(datas, &storagetransfer.HttpData{
			ListUrl: data["list_url"].(string),
		})
	}
	return datas
}

func flattenHttpData(httpDatas []*storagetransfer.HttpData) []map[string]interface{} {
	datasSchema := make([]map[string]interface{}, 0, len(httpDatas))
	for _, data := range httpDatas {
		datasSchema = append(datasSchema, map[string]interface{}{
			"list_url": data.ListUrl,
		})
	}
	return datasSchema
}

func expandObjectConditions(conditions []interface{}) []*storagetransfer.ObjectConditions {
	datas := make([]*storagetransfer.ObjectConditions, 0, len(conditions))
	for _, raw := range conditions {
		data := raw.(map[string]interface{})
		datas = append(datas, &storagetransfer.ObjectConditions{
			ExcludePrefixes:                     convertStringArr(data["exclude_prefixes"].([]interface{})),
			IncludePrefixes:                     convertStringArr(data["include_prefixes"].([]interface{})),
			MaxTimeElapsedSinceLastModification: data["max_time_elapsed_since_last_modification"].(string),
			MinTimeElapsedSinceLastModification: data["min_time_elapsed_since_last_modification"].(string),
		})
	}
	return datas
}

func flattenObjectConditions(conditions []*storagetransfer.ObjectConditions) []map[string]interface{} {
	datasSchema := make([]map[string]interface{}, 0, len(conditions))
	for _, data := range conditions {
		datasSchema = append(datasSchema, map[string]interface{}{
			"exclude_prefixes":                         data.ExcludePrefixes,
			"include_prefixes":                         data.IncludePrefixes,
			"max_time_elapsed_since_last_modification": data.MaxTimeElapsedSinceLastModification,
			"min_time_elapsed_since_last_modification": data.MinTimeElapsedSinceLastModification,
		})
	}
	return datasSchema
}

func expandTransferOptions(options []interface{}) []*storagetransfer.TransferOptions {
	datas := make([]*storagetransfer.TransferOptions, 0, len(options))
	for _, raw := range options {
		data := raw.(map[string]interface{})
		datas = append(datas, &storagetransfer.TransferOptions{
			DeleteObjectsFromSourceAfterTransfer:  data["delete_objects_from_source_after_transfer"].(bool),
			DeleteObjectsUniqueInSink:             data["delete_objects_unique_in_sink"].(bool),
			OverwriteObjectsAlreadyExistingInSink: data["overwrite_objects_already_existing_in_sink"].(bool),
		})
	}
	return datas
}

func flattenTransferOptions(options []*storagetransfer.TransferOptions) []map[string]interface{} {
	datasSchema := make([]map[string]interface{}, 0, len(options))
	for _, data := range options {
		datasSchema = append(datasSchema, map[string]interface{}{
			"delete_objects_from_source_after_transfer":  data.DeleteObjectsFromSourceAfterTransfer,
			"delete_objects_unique_in_sink":              data.DeleteObjectsUniqueInSink,
			"overwrite_objects_already_existing_in_sink": data.OverwriteObjectsAlreadyExistingInSink,
		})
	}
	return datasSchema
}

func expandTransferSpecs(transferSpecs []interface{}) []*storagetransfer.TransferSpec {
	specs := make([]*storagetransfer.TransferSpec, 0, len(transferSpecs))
	for _, raw := range transferSpecs {
		spec := raw.(map[string]interface{})

		transferSpec := &storagetransfer.TransferSpec{
			GcsDataSink: expandGcsData(spec["gcs_data_sink"].([]interface{}))[0],
		}

		if v, ok := spec["object_conditions"]; ok && len(v.([]interface{})) > 0 {
			transferSpec.ObjectConditions = expandObjectConditions(v.([]interface{}))[0]
		}
		if v, ok := spec["transfer_options"]; ok && len(v.([]interface{})) > 0 {
			transferSpec.TransferOptions = expandTransferOptions(v.([]interface{}))[0]
		}

		if v, ok := spec["gcs_data_source"]; ok && len(v.([]interface{})) > 0 {
			transferSpec.GcsDataSource = expandGcsData(v.([]interface{}))[0]
		} else if v, ok := spec["aws_s3_data_source"]; ok && len(v.([]interface{})) > 0 {
			transferSpec.AwsS3DataSource = expandAwsS3Data(v.([]interface{}))[0]
		} else if v, ok := spec["http_data_source"]; ok && len(v.([]interface{})) > 0 {
			transferSpec.HttpDataSource = expandHttpData(v.([]interface{}))[0]
		}

		specs = append(specs, transferSpec)
	}
	return specs
}

func flattenTransferSpecs(transferSpecs []*storagetransfer.TransferSpec) []map[string][]map[string]interface{} {
	transferSpecsSchema := make([]map[string][]map[string]interface{}, 0, len(transferSpecs))
	for _, transferSpec := range transferSpecs {
		schema := map[string][]map[string]interface{}{
			"gcs_data_sink": flattenGcsData([]*storagetransfer.GcsData{transferSpec.GcsDataSink}),
		}

		if transferSpec.ObjectConditions != nil {
			schema["object_conditions"] = flattenObjectConditions([]*storagetransfer.ObjectConditions{transferSpec.ObjectConditions})
		}
		if transferSpec.TransferOptions != nil {
			schema["transfer_options"] = flattenTransferOptions([]*storagetransfer.TransferOptions{transferSpec.TransferOptions})
		}
		if transferSpec.GcsDataSource != nil {
			schema["gcs_data_source"] = flattenGcsData([]*storagetransfer.GcsData{transferSpec.GcsDataSource})
		} else if transferSpec.AwsS3DataSource != nil {
			schema["aws_s3_data_source"] = flattenAwsS3Data([]*storagetransfer.AwsS3Data{transferSpec.AwsS3DataSource})
		} else if transferSpec.HttpDataSource != nil {
			schema["http_data_source"] = flattenHttpData([]*storagetransfer.HttpData{transferSpec.HttpDataSource})
		}

		transferSpecsSchema = append(transferSpecsSchema, schema)
	}
	return transferSpecsSchema
}

func validateDuration() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if _, err := time.ParseDuration(v); err != nil {
			es = append(es, fmt.Errorf("expected %s to be a duration, but parsing gave an error: %s", k, err.Error()))
			return
		}

		return
	}
}
