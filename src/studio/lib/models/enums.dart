/// 状态枚举 —— 统一 `qtcloud-course` enums.dart（LMS 部分）与 `qtclass` session.dart。
library;

/// 班级状态。
enum ClassStatus {
  preparing,
  active,
  ended;

  String get label {
    switch (this) {
      case ClassStatus.preparing:
        return '筹备中';
      case ClassStatus.active:
        return '进行中';
      case ClassStatus.ended:
        return '已结束';
    }
  }

  static ClassStatus fromString(String value) {
    switch (value) {
      case 'active':
        return ClassStatus.active;
      case 'ended':
        return ClassStatus.ended;
      default:
        return ClassStatus.preparing;
    }
  }
}

/// 考核类型。
enum AssessmentType {
  homework,
  exam;

  String get label {
    switch (this) {
      case AssessmentType.homework:
        return '作业';
      case AssessmentType.exam:
        return '考试';
    }
  }

  static AssessmentType fromString(String value) {
    switch (value) {
      case 'exam':
        return AssessmentType.exam;
      default:
        return AssessmentType.homework;
    }
  }
}

/// 提交状态。
enum SubmissionStatus {
  submitted,
  late,
  resubmitted;

  String get label {
    switch (this) {
      case SubmissionStatus.submitted:
        return '已提交';
      case SubmissionStatus.late:
        return '迟交';
      case SubmissionStatus.resubmitted:
        return '已重交';
    }
  }

  static SubmissionStatus fromString(String value) {
    switch (value) {
      case 'late':
        return SubmissionStatus.late;
      case 'resubmitted':
        return SubmissionStatus.resubmitted;
      default:
        return SubmissionStatus.submitted;
    }
  }
}

/// 课次状态（`qtclass` SessionStatus）。
enum SessionStatus {
  upcoming,
  inProgress,
  completed;

  String get label {
    switch (this) {
      case SessionStatus.upcoming:
        return '即将开始';
      case SessionStatus.inProgress:
        return '进行中';
      case SessionStatus.completed:
        return '已结束';
    }
  }

  static SessionStatus fromString(String value) {
    switch (value) {
      case 'inProgress':
        return SessionStatus.inProgress;
      case 'completed':
        return SessionStatus.completed;
      default:
        return SessionStatus.upcoming;
    }
  }
}

/// 考勤状态（`qtclass` AttendanceStatus）。
enum AttendanceStatus {
  unknown,
  present,
  late,
  absent;

  String get label {
    switch (this) {
      case AttendanceStatus.unknown:
        return '未记';
      case AttendanceStatus.present:
        return '正常';
      case AttendanceStatus.late:
        return '迟到';
      case AttendanceStatus.absent:
        return '缺勤';
    }
  }

  static AttendanceStatus fromString(String value) {
    switch (value) {
      case 'present':
        return AttendanceStatus.present;
      case 'late':
        return AttendanceStatus.late;
      case 'absent':
        return AttendanceStatus.absent;
      default:
        return AttendanceStatus.unknown;
    }
  }
}

/// 学员权益计划。
enum PlanType {
  free,
  paid,
  vip;

  String get label {
    switch (this) {
      case PlanType.free:
        return '免费';
      case PlanType.paid:
        return '付费';
      case PlanType.vip:
        return 'VIP';
    }
  }

  static PlanType fromString(String value) {
    switch (value) {
      case 'paid':
        return PlanType.paid;
      case 'vip':
        return PlanType.vip;
      default:
        return PlanType.free;
    }
  }
}
